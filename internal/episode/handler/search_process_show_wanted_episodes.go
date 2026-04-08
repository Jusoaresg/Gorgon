package episode

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/config"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	"github.com/jusoaresg/gorgon/internal/episode/service"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Search Process Show Wanted Episodes
// @Description Search and process show wanted episodes ( Automatic Show Search )
// @Tags Database/Show/Episode
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/episode/search/all [post]
func SearchProcessShowWantedEpisodes(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Search Process Show Wanted Episodes", slog.String("endpoint", "/database/show/episode/search/all"), slog.String("method", c.Request().Method))

	db := config.GetSQLite()

	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 500, "Failed to bind body")
		return err
	}

	episodeRepo := episodeRepository.NewEpisodeRepository(db)
	showRepo := showRepository.NewShowRepository(db)

	episodeSearchService := service.NewEpisodeSearchService(
		db,
		logger,
		&service.EpisodeSearcher{},
		&service.EpisodeDownloader{},
		episodeRepo,
		showRepo,
	)

	episodeSearchService.ProcessShowWantedEpisodes(int(request.Id))

	schemas.SendSuccess(c, "Process Show Wanted Episodes", request)
	return nil
}
