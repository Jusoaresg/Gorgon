package episode

import (
	"gorgon/config"
	event "gorgon/internal/db/events/episode"
	"gorgon/internal/db/repository"
	"gorgon/internal/db/schema/episode"
	"gorgon/pkg/schemas"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Change Episode Status
// @Description Change Episode Tracking Status
// @Tags Database/Show/Episode
// @Produce json
// @Param request body episode.ChangeEpisodeTrackingRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/episode/status [post]
func ChangeEpisodeStatus(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Change Episode Status", slog.String("endpoint", "/database/show/episode/:id/status"), slog.String("method", "post"))

	var request episode.ChangeEpisodeTrackingRequest

	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 500, "Failed to bind body")
		return err
	}

	episodeRepository := repository.NewEpisodeRepository(config.GetSQLite())

	ep, err := episodeRepository.GetByID(int64(request.EpisodeId))
	if err != nil {
		logger.Error("Error while fetching episode from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching episode")
		return err
	}

	ep.Tracking = request.Tracking

	if err := episodeRepository.Update(ep); err != nil {
		schemas.SendError(c, 500, "Failed to update episode tracking status")
		return err
	}

	event.EmitEpisodeTrackingUpdatedEvent(ep.ID, ep.Tracking)

	logger.Info("Successfully updated episode", slog.Any("Episode", ep))
	schemas.SendSucess(c, "Get Show", ep)
	return nil
}
