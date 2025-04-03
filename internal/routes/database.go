package routes

import (
	"gorgon/config"
	"gorgon/internal/db/handler/anime"
	"gorgon/internal/db/handler/indexer"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupDatabaseRouter(v1 *gin.RouterGroup) {
	logger := config.GetLogger()

	listGroup := v1.Group("database")
	{
		animeGroup := listGroup.Group("anime")
		{
			animeGroup.POST("", anime.AddAnimeToList)
			logger.Info("POST route added to /api/v1/database/anime")

			animeGroup.GET("", anime.ListAnimes)
			logger.Info("GET route added to /api/v1/database/anime")

			animeGroup.DELETE("", anime.DeleteAnime)
			logger.Info("DELETE route added to /api/v1/database/anime")
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
