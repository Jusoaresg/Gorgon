package show

import (
	"github.com/jusoaresg/gorgon/config"
	episodeRepo "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonRepo "github.com/jusoaresg/gorgon/internal/season/repository"
	showRepo "github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/internal/show/service"
	showAliasRepo "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
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
	db := config.GetSQLite()

	showRepo := showRepo.NewShowRepository(db)
	showAliasRepo := showAliasRepo.NewShowAliasesRepository(db)
	episodeRepo := episodeRepo.NewEpisodeRepository(db)
	seasonRepo := seasonRepo.NewSeasonRepository(db)

	aggregatorService := service.NewShowAggregatorService(showRepo, &showAliasRepo, episodeRepo, seasonRepo)

	shows, err := aggregatorService.ListFullShows()
	if err != nil {
		logger.Error("Error while fetching shows from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching shows")
		return err
	}

	logger.Info("Successfully fetched full shows", slog.Int("count", len(shows)))
	schemas.SendSuccess(c, "List Full Shows", shows)
	return nil
}
