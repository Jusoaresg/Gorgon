package routes

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/qbittorrent/handler"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func SetupQbittorrentRouter(r *echo.Group) {
	logger := config.GetLogger()

	qbitorrentRouter := r.Group("qbittorrent/")
	{
		qbitorrentRouter.POST("add", handler.AddNewTorrent)
		logger.Info("POST route added to /api/v1/qbittorrent/add")

		qbitorrentRouter.POST("add/episode", handler.AddEpisodeTorrent)
		logger.Info("POST route added to /api/v1/qbittorrent/add/episode")

		qbitorrentRouter.GET("info", handler.CheckTorrentInfo)
		logger.Info("GET route added to /api/v1/qbittorrent/info")

		qbitorrentRouter.GET("info/hash", handler.CheckTorrentInfoByHash)
		logger.Info("GET route added to /api/v1/qbittorrent/info/hash")
	}
	logger.Info("QBittorrent routes added successfully", slog.String("endpoint", "/api/v1/qbittorrent"))
}
