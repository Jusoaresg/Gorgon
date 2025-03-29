package routes

import (
	"gorgon/pkg/handler/db"

	"github.com/gin-gonic/gin"
)

func SetupDatabaseRouter(v1 *gin.RouterGroup) {

	listGroup := v1.Group("database")
	{
		animeGroup := listGroup.Group("anime")
		{
			animeGroup.POST("", handler.AddAnimeToList)
			animeGroup.GET("", handler.ListAnimes)
			animeGroup.DELETE("", handler.DeleteAnime)
		}

		indexerGroup := listGroup.Group("indexer")
		{
			indexerGroup.GET(":id", handler.GetIndexer)
		}
	}
}
