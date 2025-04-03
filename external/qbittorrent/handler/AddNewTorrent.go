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

// @Summary Add Torrent
// @Description Add New Torrent
// @Tags QBittorrent
// @Accept json
// @Produce json
// @Param request body schema.AddNewTorrentRequest true "Request Body"
// @Success 200 {object} schema.AddNewTorrentRequest
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /qbittorrent/add [post]
func AddNewTorrent(c *gin.Context) {
	logger := config.GetLogger()
	logger.Info("Received request to AddNewTorrent", slog.String("endpoint", "/api/v1/qbittorrent/add"), slog.String("method", "POST"))

	var request schema.AddNewTorrentRequest
	if err := c.BindJSON(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("endpoint", "/api/v1/qbittorrent/add"), slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to bind request body")
		return
	}

	if request.MagneticUrl == "" {
		logger.Warn("Magnetic Url is required but received empty", slog.String("endpoint", "/api/v1/qbittorrent/add"))
		schemas.SendError(c, 400, "Magnetic Url is required")
		return
	}

	torrentService, err := service.NewQBittorrentService(logger)
	if err != nil {
		logger.Error("Error while initializing qbittorrent service", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while initializing qbittorrent service: %s", err.Error()))
		return
	}

	if err := torrentService.AddTorrent(request.MagneticUrl); err != nil {
		logger.Error("Error while adding new torrent", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while adding torrent: %s", err.Error()))
	}

	schemas.SendSucess(c, "AddNewTorrent", request)
}
