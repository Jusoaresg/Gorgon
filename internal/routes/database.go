package routes

import (
	"gorgon/config"
	"gorgon/internal/db/handler/indexer"
	"gorgon/internal/db/handler/show"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func SetupDatabaseRouter(v1 *echo.Group) {
	logger := config.GetLogger()

	listGroup := v1.Group("database/")
	{
		showGroup := listGroup.Group("show")
		{
			showGroup.POST("", show.AddShowToList)
			logger.Info("POST route added to /api/v1/database/show")

			showGroup.GET("", show.ListShows)
			logger.Info("GET route added to /api/v1/database/show")

			showGroup.DELETE("", show.DeleteShow)
			logger.Info("DELETE route added to /api/v1/database/show")
		}
		// animeGroup := listGroup.Group("anime")
		// {
		// 	animeGroup.POST("", anime.AddAnimeToList)
		// 	logger.Info("POST route added to /api/v1/database/anime")
		//
		//
		// }

		indexerGroup := listGroup.Group("indexer")
		{
			indexerGroup.GET(":id", indexer.GetIndexer)
			logger.Info("GET route added to /api/v1/database/indexer/:id")
			indexerGroup.GET("", indexer.ListIndexers)
			logger.Info("POST route added to /api/v1/database/indexer")

			indexerGroup.POST("", indexer.AddIndexer)
			logger.Info("POST route added to /api/v1/database/indexer")

			indexerGroup.DELETE("", indexer.DeleteIndexer)
			logger.Info("DELETE route added to /api/v1/database/indexer")
		}
	}
	logger.Info("Database routes added successfully", slog.String("endpoint", "/api/v1/database"))
}
