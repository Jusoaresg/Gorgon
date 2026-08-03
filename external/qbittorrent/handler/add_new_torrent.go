package handler

import (
	"errors"
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	"github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"log/slog"

	"github.com/labstack/echo/v4"
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
func AddNewTorrent(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("received request to AddNewTorrent", slog.String("endpoint", "/api/v1/qbittorrent/add"), slog.String("method", "POST"))

	var request schema.AddNewTorrentRequest
	if err := c.Bind(&request); err != nil {
		logger.Error("failed to bind request body", slog.String("endpoint", "/api/v1/qbittorrent/add"), slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to bind request body")
		return err
	}

	if request.MagneticUrl == "" {
		logger.Warn("magnetic Url is required but received empty", slog.String("endpoint", "/api/v1/qbittorrent/add"))
		schemas.SendError(c, 400, "Magnetic Url is required")
		return errors.New("Magnetic Url is required")
	}

	torrentService, err := service.NewQBittorrentService(logger)
	if err != nil {
		logger.Error("error while initializing qbittorrent service", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while initializing qbittorrent service: %s", err.Error()))
		return err
	}

	if err := torrentService.AddTorrent(request.MagneticUrl); err != nil {
		logger.Error("error while adding new torrent", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while adding torrent: %s", err.Error()))
	}

	schemas.SendSuccess(c, "AddNewTorrent", request)
	return nil
}
