package handler

import (
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	"github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Check Torrent Info By Hash
// @Description Check torrent Info By Hash
// @Tags QBittorrent
// @Accept json
// @Produce json
// @Param request query schema.CheckTorrentByHashRequest true "Request Query Parameters"
// @Success 200 {object} schema.CheckTorrentResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /qbittorrent/info/hash [get]
func CheckTorrentInfoByHash(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Check Torrent Info By Hash", slog.String("endpoint", "/api/v1/qbittorrent/info/hash"))

	status := c.QueryParam("status")
	hash := c.QueryParam("hash")
	if status == "" {
		logger.Error("Missing required 'status' query parameter")
		schemas.SendError(c, 400, "Missing 'status' query parameter")
		return fmt.Errorf("Missing 'status' query parameter")
	}

	torrentService, err := service.NewQBittorrentService(logger)
	if err != nil {
		logger.Error("Failed to create Torrent Service", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to create Torrent Service")
		return err
	}

	var response []schema.CheckTorrentResponse
	if err := torrentService.CheckTorrentsWithHash(status, hash, &response); err != nil {
		logger.Error("Error while getting torrent info by hash", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error get torrent info by hash: %s", err.Error()))
	}

	logger.Info("Check Torrent info by hash request successfully", slog.Any("response", response))
	schemas.SendSucess(c, "CheckTorrentInfoByHash", response)
	return nil
}
