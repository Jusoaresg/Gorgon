package routes

import (
	"gorgon/config"
	"gorgon/external/prowlarr/handler"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func SetupProwlarrRouter(v1 *echo.Group) {
	logger := config.GetLogger()

	prowlarrGroup := v1.Group("prowlarr/")
	{
		indexerGroup := prowlarrGroup.Group("indexer")
		{
			indexerGroup.GET("", handler.GetIndexer)
			logger.Info("GET route added to /api/v1/prowlarr/indexer")
		}

		searchGroup := prowlarrGroup.Group("search")
		{
			searchGroup.POST("", handler.SearchAnimes)
			logger.Info("POSTS route added to /api/v1/prowlarr/search")
		}
	}
	logger.Info("Prowlarr routes added successfully", slog.String("endpoint", "/api/v1/prowlarr"))
}
