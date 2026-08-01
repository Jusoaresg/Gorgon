package handler

import (
	"errors"
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	"github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/episode/events"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	"github.com/jusoaresg/gorgon/internal/episode/repository"
	episodeTorrentModel "github.com/jusoaresg/gorgon/internal/episode_torrent/model"
	episodeTorrentRepository "github.com/jusoaresg/gorgon/internal/episode_torrent/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"log/slog"

	"github.com/jmoiron/sqlx"
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
// @Router /qbittorrent/add/episode [post]
func AddEpisodeTorrent(c echo.Context) error {
	return addEpisodeTorrentHandler(c, config.GetSQLite())
}

func addEpisodeTorrentHandler(c echo.Context, db *sqlx.DB) error {
	logger := config.GetLogger()
	logger.Info("Received request to AddNewTorrent", slog.String("endpoint", "/api/v1/qbittorrent/add"), slog.String("method", "POST"))

	var request schema.AddEpisodeTorrentRequest
	if err := c.Bind(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("endpoint", "/api/v1/qbittorrent/add"), slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to bind request body")
		return err
	}
	epRepo := repository.NewEpisodeRepository(db)
	episodeTorrentRepo := episodeTorrentRepository.NewEpisodeTorrentRepository(db)
	ep, err := epRepo.GetByID(request.EpisodeID)
	if err != nil {
		logger.Error("Error while fetching episode id", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while adding episode torrent: %s", err.Error()))
		return err
	}

	if request.MagneticUrl == "" {
		logger.Warn("Magnetic Url is required but received empty", slog.String("endpoint", "/api/v1/qbittorrent/add"))
		schemas.SendError(c, 400, "Magnetic Url is required")
		return errors.New("Magnetic Url is required")
	}

	torrentService, err := service.NewQBittorrentService(logger)
	if err != nil {
		logger.Error("Error while initializing qbittorrent service", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while initializing qbittorrent service: %s", err.Error()))
		return err
	}

	if err := torrentService.AddTorrent(request.MagneticUrl); err != nil {
		logger.Error("Error while adding new torrent", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while adding torrent: %s", err.Error()))
	}

	ep.Tracking = model.TrackingSnatched

	episodeTorrent := episodeTorrentModel.FromAddTorrentRequest(ep.ID, request)
	if _, err := episodeTorrentRepo.Upsert(episodeTorrent); err != nil {
		logger.Error("Error while saving episode torrent", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while adding episode torrent: %s", err.Error()))
		return err
	}

	if err := epRepo.Update(ep); err != nil {
		logger.Error("Error while updating episode tracking when adding new torrent", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while adding episode torrent: %s", err.Error()))
		return err
	}

	episode.EmitEpisodeTrackingUpdatedEvent(ep.ID, ep.Tracking)

	schemas.SendSuccess(c, "AddNewTorrent", request)
	return nil
}
