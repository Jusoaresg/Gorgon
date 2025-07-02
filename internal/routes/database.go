package routes

import (
	"github.com/jusoaresg/gorgon/config"
	episodeHandler "github.com/jusoaresg/gorgon/internal/episode/handler"
	indexerHandler "github.com/jusoaresg/gorgon/internal/indexer/handler"
	seasonHandler "github.com/jusoaresg/gorgon/internal/season/handler"
	showHandler "github.com/jusoaresg/gorgon/internal/show/handler"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func SetupDatabaseRouter(v1 *echo.Group) {
	logger := config.GetLogger()

	listGroup := v1.Group("database/")
	{
		showGroup := listGroup.Group("show")
		{
			showGroup.POST("", showHandler.AddShowToList)
			logger.Info("POST route added to /api/v1/database/show")

			showGroup.GET("", showHandler.ListShows)
			logger.Info("GET route added to /api/v1/database/show")

			showGroup.GET("/full", showHandler.ListFullShows)
			logger.Info("GET route added to /api/v1/database/show/full")

			showGroup.GET("/:id", showHandler.GetShow)
			logger.Info("GET route added to /api/v1/database/show/:id")

			showGroup.DELETE("", showHandler.DeleteShow)
			logger.Info("DELETE route added to /api/v1/database/show")

			showGroup.POST("/update-info", showHandler.UpdateShowInfo)
			logger.Info("POST route added to /api/v1/database/show/update-info")

			episodeGroup := showGroup.Group("/episode")
			{
				episodeGroup.POST("/status", episodeHandler.ChangeEpisodeStatus)
				logger.Info("POST route added to /api/v1/database/show/episode/status")

				episodeGroup.GET("/:id", episodeHandler.GetShowEpisodes)
				logger.Info("GET route added to /api/v1/database/show/episode/:id")

				episodeGroup.DELETE("/:id", episodeHandler.DeleteDownloadedEpisode)
				logger.Info("DELETE route added to /api/v1/database/show/episode/:id")
			}

			seasonsGroup := showGroup.Group("/season")
			{
				seasonsGroup.GET("/:id", seasonHandler.GetShowSeasons)
				logger.Info("GET route added to /api/v1/database/show/season/:id")
			}
		}

		indexerGroup := listGroup.Group("indexer")
		{
			indexerGroup.GET(":id", indexerHandler.GetIndexer)
			logger.Info("GET route added to /api/v1/database/indexer/:id")
			indexerGroup.GET("", indexerHandler.ListIndexers)
			logger.Info("POST route added to /api/v1/database/indexer")

			indexerGroup.POST("", indexerHandler.AddIndexer)
			logger.Info("POST route added to /api/v1/database/indexer")

			indexerGroup.DELETE("", indexerHandler.DeleteIndexer)
			logger.Info("DELETE route added to /api/v1/database/indexer")
		}
	}
	logger.Info("Database routes added successfully", slog.String("endpoint", "/api/v1/database"))
}
