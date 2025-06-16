package handler

import (
	"gorgon/config"
	"gorgon/external/tvmaze/schema"
	"gorgon/external/tvmaze/service"
	"gorgon/internal/db/repository"
	"gorgon/pkg/schemas"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Seach Show By Name
// @Description Search TvMaze By Show Name
// @Tags TvMaze/Search
// @Accept json
// @Produce json
// @Param request body schemas.NameRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /tvmaze/search/name [post]
func SearchShowByName(c echo.Context) error {
	logger := config.GetLogger()

	var request schemas.NameRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	tvMazeService := service.NewTvMazeSearchService(logger)

	response, err := tvMazeService.SearchByName(request.Name)
	if err != nil {
		logger.Error("Error while searching for name", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while searching by name")
		return err
	}

	showRepo := repository.NewShowRepository()
	existingShows, err := showRepo.List()
	if err != nil {
		return err
	}

	existingsMap := make(map[int64]bool)
	for _, s := range existingShows {
		existingsMap[s.TvMazeID] = true
	}

	var enriched []schema.SearchResult
	for _, r := range *response {
		enriched = append(enriched, schema.SearchResult{
			Show:    r.Show,
			IsAdded: existingsMap[r.Show.TvMazeID],
		})
	}

	schemas.SendSucess(c, "SearchShowByName", &enriched)
	return nil
}
