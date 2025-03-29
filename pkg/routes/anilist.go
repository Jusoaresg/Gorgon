package routes

import (
	"gorgon/pkg/handler/anilist"

	"github.com/gin-gonic/gin"
)

func SetupAnilistRouter(v1 *gin.RouterGroup) {

	anilistGroup := v1.Group("anilist")
	{
		anilistGroup.POST("", anilist.FindAnimeIdAnilist)

		anilistGroup.POST("get-info", anilist.GetAnimeInfoById)

	}
}
