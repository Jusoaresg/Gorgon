package handler

import (
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/external/prowlarr/service"
	"github.com/jusoaresg/gorgon/internal/db/repository"
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

	episodeRepo := repository.NewEpisodeRepository(db)
	showRepo := repository.NewShowRepository(db)

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

	title := utils.NormalizeTitle(show.Name)
	query := fmt.Sprintf("%s {Season:%d}{Episode:%d}", title, episode.Season, episode.Number)

	searchKey := schema.SearchByTypeRequest{
		Query: query,
		Type:  "tvsearch",
	}

	searchService := service.NewProwlarrSearchService(logger)

	var response []schema.SearchResponse
	if err := searchService.SearchByType(&searchKey, &response); err != nil {
		logger.Error("Error while search show episodes from prowlarr", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while searching show episodes prowlarr: %s", err.Error()))
		return err
	}

	logger.Info("Search Shows Episodes request successfully", slog.Any("response", response))
	schemas.SendSuccess(c, "Search Shows", response)
	return nil
}
