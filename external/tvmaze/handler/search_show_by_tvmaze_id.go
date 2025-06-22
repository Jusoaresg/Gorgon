package handler

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/tvmaze/service"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Search Show By TvMaze ID
// @Description Search TvMaze By TvMaze ID
// @Tags TvMaze/Search
// @Accept json
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /tvmaze/search/tvmaze [post]
func SearchShowByTvMazeId(c echo.Context) error {
	logger := config.GetLogger()

	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	tvMazeService := service.NewTvMazeSearchService(logger)
	response, err := tvMazeService.SearchByTvMazeId(request.Id)
	if err != nil {
		logger.Error("Error while searching for name", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while searching by name")
		return err
	}

	schemas.SendSuccess(c, "Search Show By TvMaze Id", &response)
	return nil
}
