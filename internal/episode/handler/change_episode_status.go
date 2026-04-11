package episode

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/qbittorrent/service"
	event "github.com/jusoaresg/gorgon/internal/episode/events"
	"github.com/jusoaresg/gorgon/internal/episode/repository"
	"github.com/jusoaresg/gorgon/internal/episode/schema"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/jusoaresg/gorgon/utils"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Change Episodes Status
// @Description Change Episodes Tracking Status
// @Tags Database/Show/Episode
// @Produce json
// @Param request body schema.ChangeEpisodeTrackingRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/episode/status [post]
func ChangeEpisodeStatus(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Change Episode Status", slog.String("endpoint", "/database/show/episode/status"), slog.String("method", "post"))

	var request schema.ChangeEpisodeTrackingRequest

	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 500, "Failed to bind body")
		return err
	}

	episodeRepository := repository.NewEpisodeRepository(config.GetSQLite())
	qbittorrentService, err := service.NewQBittorrentService(logger)

	episodes, err := episodeRepository.GetAllByID(utils.ToInt64Slice(request.EpisodeIds)...)
	if err != nil {
		logger.Error("Error while fetching episode from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching episode")
		return err
	}

	for _, episode := range episodes {

		if episode.TorrentHash != "" {
			err := qbittorrentService.DeleteTorrent(episode.TorrentHash, true)
			if err != nil {
				logger.Error(
					"Error while deleting torrent",
					slog.String("error", err.Error()),
					slog.Int64("episode_id", episode.ID),
				)
			}
			episode.TorrentHash = ""
		}
		episode.Tracking = request.Tracking

		if err := episodeRepository.Update(episode); err != nil {
			schemas.SendError(c, 500, "Failed to update episode tracking status")
			return err
		}

		event.EmitEpisodeTrackingUpdatedEvent(episode.ID, episode.Tracking)
	}

	logger.Info("Successfully updated episodes", slog.Any("Episodes", episodes))
	c.Response().Header().Add("HX-Refresh", "true")
	schemas.SendSuccess(c, "Change episodes tracking", episodes)
	return nil
}
