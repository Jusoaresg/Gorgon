package handler

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/qbittorrent/schema"
	"gorgon/external/qbittorrent/service"
	"gorgon/pkg/schemas"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Check Torrent info
// @Description Check torrent info
// @Tags QBittorrent
// @Accept json
// @Produce json
// @Param request body schema.CheckTorrentRequest true "Request Body"
// @Success 200 {object} schema.CheckTorrentResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /qbittorrent/info [get]
func CheckTorrentInfo(c *gin.Context) {
	logger := config.GetLogger()
	logger.Info("Received request to Check Torrent Info", slog.String("endpoint", "/api/v1/qbittorrent/info"))

	var request schema.CheckTorrentRequest
	if err := c.BindJSON(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to bind request body")
		return
	}

	torrentService, err := service.NewQBittorrentService(logger)
	if err != nil {
		logger.Error("Failed to create Torrent Service", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to create Torrent Service")
		return
	}

	var response []schema.CheckTorrentResponse
	if err := torrentService.CheckTorrents("all", &response); err != nil {
		logger.Error("Error while getting torrent info", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error get torrent info: %s", err.Error()))
	}

	logger.Info("Check Torrent info request successfully", slog.Any("response", response))
	schemas.SendSucess(c, "CheckTorrentInfo", response)
}
