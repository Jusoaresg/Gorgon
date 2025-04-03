package routes

import (
	"gorgon/config"
	"gorgon/external/tvmaze/handler"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func SetupTvMazeRouter(v1 *echo.Group) {
	logger := config.GetLogger()

	tvmazeGroup := v1.Group("tvmaze/")
	{
		searchGroup := tvmazeGroup.Group("search/")
		{
			if err := searchGroup.POST("name", handler.SearchShowByName); err != nil {
				return
			}

			logger.Info("POST route added to /api/v1/tvmaze/search/name")

			searchGroup.POST("tvmaze", handler.SearchShowByTvMazeId)
			logger.Info("POST route added to /api/v1/tvmaze/search/tvmaze")
		}

	}
	logger.Info("TvMaze routes added succesfully", slog.String("endpoint", "/api/v1/tvmaze/"))
}
