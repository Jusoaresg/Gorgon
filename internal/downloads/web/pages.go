package web

import (
	"log/slog"
	"net/http"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	qbittorrentService "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/downloads/service"
	"github.com/jusoaresg/gorgon/views"
	"github.com/labstack/echo/v4"
)

type DownloadsData struct {
	Items        []service.DownloadItem
	ErrorMessage string
}

func (h *Handler) DownloadsRoute(c echo.Context) error {
	items, err := h.fetchDownloadItems()
	errorMessage := ""
	if err != nil {
		errorMessage = "Unable to reach the torrent client. Configure qBittorrent in Settings."
	}

	return views.Render(c, views.View{
		Layout:    "layout",
		Default:   "downloads",
		Templates: map[string]string{"download-items": "download-items"},
		Data:      DownloadsData{Items: items, ErrorMessage: errorMessage},
		Styles:    []string{"downloads.css"},
	})
}

func (h *Handler) DownloadsItemsHTMX(c echo.Context) error {
	items, err := h.fetchDownloadItems()
	if err != nil {
		return c.Render(http.StatusOK, "download-items", views.PageData{Data: DownloadsData{Items: []service.DownloadItem{}}})
	}

	return c.Render(http.StatusOK, "download-items", views.PageData{Data: DownloadsData{Items: items}})
}

func (h *Handler) fetchDownloadItems() ([]service.DownloadItem, error) {
	logger := config.GetLogger()

	torrentService, err := qbittorrentService.NewQBittorrentService(logger)
	if err != nil {
		logger.Warn("failed to create qbittorrent service for downloads page", slog.String("error", err.Error()))
		return []service.DownloadItem{}, err
	}

	var torrents []schema.CheckTorrentResponse
	if err := torrentService.CheckTorrents("all", &torrents); err != nil {
		logger.Warn("failed to fetch torrents for downloads page", slog.String("error", err.Error()))
		return []service.DownloadItem{}, err
	}

	return h.Service.BuildDownloads(torrents)
}
