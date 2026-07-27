package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jusoaresg/gorgon/config"
	prowlarrSchema "github.com/jusoaresg/gorgon/external/prowlarr/schema"
	prowlarrService "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrentService "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	tvMazeService "github.com/jusoaresg/gorgon/external/tvmaze/service"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	episodeEvents "github.com/jusoaresg/gorgon/internal/episode/events"
	seasonModel "github.com/jusoaresg/gorgon/internal/season/model"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	showSchema "github.com/jusoaresg/gorgon/internal/show/schema"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showManager "github.com/jusoaresg/gorgon/internal/show/service"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/jusoaresg/gorgon/pkg/schemas/dtos"
	"github.com/jusoaresg/gorgon/pkg/services"
	"github.com/jusoaresg/gorgon/utils"
	"github.com/labstack/echo/v4"
)

func (h *Handler) SearchShowHTMX(c echo.Context) error {
	var request schemas.NameRequest
	if err := c.Bind(&request); err != nil {
		return nil
	}

	shows, err := h.ShowManager.SearchAndEnrich(request.Name)
	if err != nil {
		return err
	}

	data := map[string]any{
		"Shows": shows,
	}

	return c.Render(http.StatusOK, "add-show-card", data)
}

func (h *Handler) AddShowHTMX(c echo.Context) error {
	id, err := strconv.ParseInt(c.FormValue("id"), 10, 64)
	if err != nil {
		return err
	}

	trackingType := c.FormValue("monitor")

	request := showSchema.AddShowToListRequest{
		Id:           int(id),
		TrackingType: trackingType,
	}

	validTrackings := map[string]bool{
		"all":    true,
		"future": true,
		"none":   true,
	}

	if !validTrackings[request.TrackingType] {
		return echo.NewHTTPError(400, "Invalid tracking type")
	}

	logger := config.GetLogger()

	tvMazeSvc := tvMazeService.NewTvMazeSearchService(logger)
	showManagerSvc := showManager.NewShowManagerService(logger, h.DB)

	idString := strconv.Itoa(request.Id)
	id64, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return err
	}

	showDto, err := tvMazeSvc.SearchByTvMazeId(id64)
	if err != nil {
		return err
	}

	episodesDto, err := showManagerSvc.GetEpisodes(showDto.TvMazeID)
	if err != nil {
		return err
	}

	services.ApplyTrackingToEpisodes(episodesDto, request.TrackingType)

	seasonsDto, err := showManagerSvc.GetSeasons(showDto.TvMazeID)
	if err != nil {
		return err
	}

	show := showDto.ToModel()

	tx, err := h.DB.Beginx()
	if err != nil {
		return err
	}

	showRepo := showRepository.NewShowRepository(h.DB)
	showID, err := showRepo.CreateTx(tx, show)
	if err != nil {
		logger.Error("Failed to add show to database", slog.String("error", err.Error()))
		return err
	}

	aliasRepo := showAliasRepository.NewShowAliasesRepository(h.DB)
	aliases := showDto.ToAliasModel()
	for _, alias := range aliases {
		alias.ShowID = showID
		alias.Source = "tvmaze"
		_, err := aliasRepo.CreateTx(tx, alias)
		if err != nil {
			logger.Error(
				"failed to create show alias",
				slog.String("error", err.Error()),
				slog.String("alias", alias.Alias),
				slog.String("country", alias.Country),
			)
			continue
		}
	}

	seasons := dtos.SeasonDtoSliceToModel(*seasonsDto, showID)
	var episodes []episodeModel.Episode

	for _, season := range seasons {
		seasonRepo := seasonRepository.NewSeasonRepository(h.DB)
		seasonID, err := seasonRepo.CreateTx(tx, season)
		if err != nil {
			logger.Error("Failed to add season to database", slog.String("error", err.Error()))
			return err
		}

		tempEpisodesDTO := *episodesDto
		for _, episode := range tempEpisodesDTO {
			if episode.Season == season.Number {
				ep := episode.ToModel(showID, seasonID)
				episodes = append(episodes, *ep)
			}
		}
	}

	for _, episode := range episodes {
		episodeRepo := episodeRepository.NewEpisodeRepository(h.DB)
		if _, err := episodeRepo.CreateTx(tx, episode); err != nil {
			logger.Error("Failed to add episode to database", slog.String("error", err.Error()))
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func (h *Handler) EditShowModalHTMX(c echo.Context) error {
	id := c.Param("id")
	return c.Render(http.StatusOK, "edit-show-modal", id)
}

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

type seasonModal struct {
	Season     seasonModel.Season
	EpisodeIds []int
}

func (h *Handler) ChangeSeasonTrackingModal(c echo.Context) error {
	seasonIdStr := c.Param("id")
	seasonIdInt, err := strconv.Atoi(seasonIdStr)
	if err != nil {
		return err
	}

	episodes, err := h.EpisodeRepo.ListBySeasonID(seasonIdInt)
	if err != nil {
		return err
	}

	season, err := h.SeasonRepo.GetById(int64(seasonIdInt))
	if err != nil {
		return err
	}

	episodesIds := make([]int, 0, len(episodes))
	for _, e := range episodes {
		episodesIds = append(episodesIds, int(e.ID))
	}

	return c.Render(http.StatusOK, "season-tracking-modal", seasonModal{
		Season:     season,
		EpisodeIds: episodesIds,
	})
}

type interactiveSearchData struct {
	EpisodeID   int64
	EpisodeName string
	Season      int
	Number      int
	AutoSearch  bool
}

type searchResultsData struct {
	Results   []prowlarrSchema.SearchResponse
	EpisodeID int64
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

	aliasRepo := showAliasRepository.NewShowAliasesRepository(h.DB)
	aliases, err := aliasRepo.ListByShowID(show.ID)
	if err != nil {
		logger.Error("error fetching show aliases", slog.String("error", err.Error()))
		return err
	}

	titles := []string{utils.NormalizeTitle(show.Name)}
	for _, alias := range aliases {
		titles = append(titles, alias.Alias)
	}

	searchService, err := prowlarrService.NewProwlarrSearchService(logger)
	if err != nil {
		logger.Error("error initializing prowlarr service", slog.String("error", err.Error()))
		return err
	}

	var allResults []prowlarrSchema.SearchResponse
	for _, title := range titles {
		query := fmt.Sprintf("%s S%02dE%02d", title, episode.Season, episode.Number)
		searchKey := prowlarrSchema.SearchByTypeRequest{
			Query: query,
			Type:  "tvsearch",
		}

		var results []prowlarrSchema.SearchResponse
		if err := searchService.SearchByType(&searchKey, &results); err != nil {
			logger.Error("error searching prowlarr",
				slog.String("query", query),
				slog.String("error", err.Error()),
			)
			continue
		}
		allResults = append(allResults, results...)
	}

	return c.Render(http.StatusOK, "search-results-list", searchResultsData{
		Results:   allResults,
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
