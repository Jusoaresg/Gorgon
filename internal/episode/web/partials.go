package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/jusoaresg/gorgon/config"
	prowlarrSchema "github.com/jusoaresg/gorgon/external/prowlarr/schema"
	prowlarrService "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrentService "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	episodeEvents "github.com/jusoaresg/gorgon/internal/episode/events"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/jusoaresg/gorgon/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ChangeEpisodeTrackingModal(c echo.Context) error {
	epIdStr := c.Param("id")
	epIdInt, err := strconv.Atoi(epIdStr)
	if err != nil {
		return err
	}

	episode, err := h.EpisodeRepo.GetByID(int64(epIdInt))
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "episode-tracking-modal", episode)
}

type interactiveSearchData struct {
	EpisodeID   int64
	EpisodeName string
	Season      int
	Number      int
	AutoSearch  bool
}

func (h *Handler) InteractiveSearchModal(c echo.Context) error {
	epIdStr := c.Param("id")
	epId, err := strconv.ParseInt(epIdStr, 10, 64)
	if err != nil {
		return err
	}

	episode, err := h.EpisodeRepo.GetByID(epId)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "interactive-search-modal", interactiveSearchData{
		EpisodeID:   episode.ID,
		EpisodeName: episode.Name,
		Season:      episode.Season,
		Number:      episode.Number,
		AutoSearch:  true,
	})
}

type aliasProgressItem struct {
	Name    string
	Encoded string
	Delay   string // " delay:3s" for low-priority aliases, empty for high-priority
}

var latinAliasRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-_!?',.:]+$`)

type searchProgressData struct {
	EpisodeID int64
	Season    int
	Number    int
	Titles    []aliasProgressItem
}

func (h *Handler) SearchEpisodeResults(c echo.Context) error {
	epIdStr := c.Param("id")
	epId, err := strconv.ParseInt(epIdStr, 10, 64)
	if err != nil {
		return err
	}

	logger := config.GetLogger()

	episode, err := h.EpisodeRepo.GetByID(epId)
	if err != nil {
		return err
	}

	show, err := h.ShowRepo.GetById(episode.ShowID)
	if err != nil {
		return err
	}

	aliases, err := h.ShowAliasesRepo.ListByShowID(show.ID)
	if err != nil {
		logger.Error("error fetching show aliases", slog.String("error", err.Error()))
		return err
	}

	titles := []aliasProgressItem{
		{Name: utils.NormalizeTitle(show.Name), Encoded: url.QueryEscape(utils.NormalizeTitle(show.Name))},
	}

	for _, alias := range aliases {
		if !latinAliasRegex.MatchString(alias.Alias) {
			continue
		}
		titles = append(titles, aliasProgressItem{
			Name:    alias.Alias,
			Encoded: url.QueryEscape(alias.Alias),
		})
	}

	for _, alias := range aliases {
		if latinAliasRegex.MatchString(alias.Alias) {
			continue
		}
		titles = append(titles, aliasProgressItem{
			Name:    alias.Alias,
			Encoded: url.QueryEscape(alias.Alias),
			Delay:   " delay:3s",
		})
	}

	return c.Render(http.StatusOK, "search-results-progress", searchProgressData{
		EpisodeID: epId,
		Season:    episode.Season,
		Number:    episode.Number,
		Titles:    titles,
	})
}

type searchAliasResultsData struct {
	Alias     string
	Results   []prowlarrSchema.SearchResponse
	EpisodeID int64
}

func (h *Handler) SearchAliasResult(c echo.Context) error {
	epIdStr := c.Param("id")
	epId, err := strconv.ParseInt(epIdStr, 10, 64)
	if err != nil {
		return err
	}

	title := c.QueryParam("t")
	if title == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	seasonStr := c.QueryParam("s")
	numberStr := c.QueryParam("n")
	season, _ := strconv.Atoi(seasonStr)
	number, _ := strconv.Atoi(numberStr)

	logger := config.GetLogger()

	prowlarrIndexerService := prowlarrService.NewProwlarrIndexerService(logger)
	var indexers []prowlarrSchema.IndexerResponse
	if err := prowlarrIndexerService.GetIndexers(&indexers); err != nil {
		logger.Error("error fetching prowlarr indexers", slog.String("error", err.Error()))
		return c.Render(http.StatusOK, "search-alias-results", searchAliasResultsData{
			Alias:     title,
			Results:   nil,
			EpisodeID: epId,
		})
	}

	var indexerIds []int
	for _, indexer := range indexers {
		if indexer.Enabled {
			indexerIds = append(indexerIds, indexer.Id)
		}
	}

	searchService, err := prowlarrService.NewInteractiveProwlarrSearchService(logger)
	if err != nil {
		logger.Error("error initializing prowlarr service", slog.String("error", err.Error()))
		return c.Render(http.StatusOK, "search-alias-results", searchAliasResultsData{
			Alias:     title,
			Results:   nil,
			EpisodeID: epId,
		})
	}

	query := fmt.Sprintf("%s S%02dE%02d", title, season, number)
	searchKey := prowlarrSchema.SearchByTypeRequest{
		Query: query,
		Type:  "tvsearch",
	}

	var results []prowlarrSchema.SearchResponse
	if err := searchService.SearchByType(&searchKey, &results, indexerIds...); err != nil {
		logger.Error("error searching prowlarr",
			slog.String("query", query),
			slog.String("error", err.Error()),
		)
	}

	return c.Render(http.StatusOK, "search-alias-results", searchAliasResultsData{
		Alias:     title,
		Results:   results,
		EpisodeID: epId,
	})
}

func (h *Handler) DownloadEpisodeTorrent(c echo.Context) error {
	epIdStr := c.Param("id")
	epId, err := strconv.ParseInt(epIdStr, 10, 64)
	if err != nil {
		schemas.SendError(c, 400, "Invalid episode ID")
		return nil
	}

	var request struct {
		Guid     string `json:"guid"`
		InfoHash string `json:"infoHash"`
	}
	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 400, "Invalid request")
		return nil
	}

	logger := config.GetLogger()

	ep, err := h.EpisodeRepo.GetByID(epId)
	if err != nil {
		logger.Error("error fetching episode", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Episode not found")
		return nil
	}

	torrentService, err := qbittorrentService.NewQBittorrentService(logger)
	if err != nil {
		logger.Error("error initializing qbittorrent service", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Torrent client not available")
		return nil
	}

	if err := torrentService.AddTorrent(request.Guid); err != nil {
		logger.Error("error adding torrent", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to add torrent")
		return nil
	}

	ep.Tracking = episodeModel.TrackingSnatched
	ep.TorrentHash = request.InfoHash

	if err := h.EpisodeRepo.Update(ep); err != nil {
		logger.Error("error updating episode tracking", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to update episode")
		return nil
	}

	episodeEvents.EmitEpisodeTrackingUpdatedEvent(ep.ID, ep.Tracking)

	logger.Info("torrent added and episode snatched",
		slog.Int64("episode_id", ep.ID),
		slog.String("hash", request.InfoHash),
	)

	schemas.SendSuccess(c, "Download started", nil)
	return nil
}
