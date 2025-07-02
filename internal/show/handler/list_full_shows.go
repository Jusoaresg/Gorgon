package show

import (
	"github.com/jusoaresg/gorgon/config"
	episodeRepo "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonRepo "github.com/jusoaresg/gorgon/internal/season/repository"
	showRepo "github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/internal/show/service"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary List Full Shows
// @Description List all shows
// @Tags Database/Show
// @Produce json
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/full [get]
func ListFullShows(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to List Full Shows", slog.String("endpoint", "/database/show/full"), slog.String("method", "get"))

	showRepo := showRepo.NewShowRepository(config.GetSQLite())
	episodeRepo := episodeRepo.NewEpisodeRepository(config.GetSQLite())
	seasonRepo := seasonRepo.NewSeasonRepository(config.GetSQLite())

	aggregatorSercice := service.NewShowAggregatorService(showRepo, episodeRepo, seasonRepo)

	shows, err := aggregatorSercice.ListFullShows()
	if err != nil {
		logger.Error("Error while fetching shows from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching shows")
		return err
	}

	logger.Info("Successfully fetched full shows", slog.Int("count", len(shows)))
	schemas.SendSuccess(c, "List Full Shows", shows)
	return nil
}
