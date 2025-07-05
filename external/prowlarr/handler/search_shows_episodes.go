package handler

import (
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/external/prowlarr/service"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/jusoaresg/gorgon/utils"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Search For Show Episodes
// @Description Search show episodes on prowlarr indexers
// @Tags Prowlarr/Search
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} []schema.SearchResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /prowlarr/search/episode [post]
func SearchShowsEpisodes(c echo.Context) error {
	return searchShowsEpisodesHandler(c, config.GetSQLite())
}

func searchShowsEpisodesHandler(c echo.Context, db *sqlx.DB) error {
	logger := config.GetLogger()
	logger.Info("Received request to Search Shows Episodes", slog.String("endpoint", "/api/v1/prowlarr/search/episode"), slog.String("method", "POST"))

	episodeRepo := episodeRepository.NewEpisodeRepository(db)
	showRepo := showRepository.NewShowRepository(db)
	showAliasRepository := showAliasRepository.NewShowAliasesRepository(db)

	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to bind request body")
		return err
	}

	episode, err := episodeRepo.GetByID(request.Id)
	if err != nil {
		return err
	}

	show, err := showRepo.GetById(episode.ShowID)
	if err != nil {
		return err
	}

	aliases, err := showAliasRepository.ListByShowID(show.ID)
	if err != nil {
		logger.Error("error fetching show aliases from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "error fetching show aliases from database")
		return err
	}
	titles := []string{
		utils.NormalizeTitle(show.Name),
	}
	for _, alias := range aliases {
		titles = append(titles, alias.Alias)
	}

	searchService, err := service.NewProwlarrSearchService(logger)
	if err != nil {
		logger.Error("error to initialize prowlarr service", slog.String("Error", err.Error()))
		return err
	}

	var response []schema.SearchResponse
	for _, title := range titles {
		query := fmt.Sprintf("%s {Season:%d}{Episode:%d}", title, episode.Season, episode.Number)
		searchKey := schema.SearchByTypeRequest{
			Query: query,
			Type:  "tvsearch",
		}

		var resp []schema.SearchResponse
		if err := searchService.SearchByType(&searchKey, &resp); err != nil {
			logger.Error(
				"error while search show episodes from prowlarr",
				slog.String("error", err.Error()),
				slog.Int64("show_id", show.ID),
				slog.String("show_alias", title),
				slog.Int64("episode_id", episode.ID),
				slog.String("episode_name", episode.Name),
			)
			schemas.SendError(c, 500, fmt.Sprintf("Error while searching show episodes prowlarr: %s", err.Error()))
			return err
		}
		response = append(response, resp...)
	}

	logger.Info("Search Shows Episodes request successfully")
	schemas.SendSuccess(c, "Search Shows", response)
	return nil
}
