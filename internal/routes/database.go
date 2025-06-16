package routes

import (
	"gorgon/config"
	"gorgon/internal/db/handler/episode"
	"gorgon/internal/db/handler/indexer"
	"gorgon/internal/db/handler/season"
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

			showGroup.GET("/:id", show.GetShow)
			logger.Info("GET route added to /api/v1/database/show/:id")

			showGroup.DELETE("", show.DeleteShow)
			logger.Info("DELETE route added to /api/v1/database/show")

			episodeGroup := showGroup.Group("/episode")
			{
				episodeGroup.POST("/status", episode.ChangeEpisodeStatus)
				logger.Info("POST route added to /api/v1/database/show/episode/status")

				episodeGroup.GET("/:id", episode.GetShowEpisodes)
				logger.Info("GET route added to /api/v1/database/show/episode/:id")
			}

			seasonsGroup := showGroup.Group("/season")
			{
				seasonsGroup.GET("/:id", season.GetShowSeasons)
				logger.Info("GET route added to /api/v1/database/show/season/:id")
			}
		}

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
