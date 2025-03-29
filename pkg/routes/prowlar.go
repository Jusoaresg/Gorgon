package routes

import (
	"gorgon/pkg/handler/prowlarr"

	"github.com/gin-gonic/gin"
)

func SetupProwlarrRouter(v1 *gin.RouterGroup) {

	prowlarrGroup := v1.Group("prowlarr")
	{
		// prowlarrGroup.POST("search", handler.AddAnimeToList)
		prowlarrGroup.GET("indexers", prowlarr.GetIndexers)
		prowlarrGroup.POST("indexers", prowlarr.AddIndexer)
		prowlarrGroup.DELETE("indexers", prowlarr.DeleteIndexer)
	}
}
