package handler

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/tvmaze/service"
	"github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Search Show By Name
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
	showRepo := repository.NewShowRepository(config.GetSQLite())
	showManager := service.NewShowManager(*tvMazeService, *showRepo, logger)

	enriched, err := showManager.SearchAndEnrich(request.Name)
	if err != nil {
		schemas.SendError(c, 500, "Error While searching by name")
		return err
	}

	schemas.SendSuccess(c, "SearchShowByName", &enriched)
	return nil
}
