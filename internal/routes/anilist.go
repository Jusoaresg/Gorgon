package routes

import (
	"gorgon/config"
	"gorgon/external/anilist/handler"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func SetupAnilistRouter(v1 *gin.RouterGroup) {
	logger := config.GetLogger()

	anilistGroup := v1.Group("anilist")
	{
		anilistGroup.POST("", handler.FindAnimeIdAnilist)
		logger.Info("POST route added to /api/v1/anilist")

		anilistGroup.POST("get-info", handler.GetAnimeInfoById)
		logger.Info("POST route added to /api/v1/anilist/get-info")
	}
	logger.Info("Anilist routes added succesfully", slog.String("endpoint", "/api/v1/anilist/"))
}
