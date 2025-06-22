package show

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/db/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary List Shows
// @Description List all shows
// @Tags Database/Show
// @Produce json
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show [get]
func ListShows(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to List Shows", slog.String("endpoint", "/database/show"), slog.String("method", "get"))

	showRepo := repository.NewShowRepository(config.GetSQLite())

	show, err := showRepo.List()
	if err != nil {
		logger.Error("Error while fetching shows from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching shows")
		return err
	}

	logger.Info("Successfully fetched shows", slog.Int("count", len(show)))
	schemas.SendSuccess(c, "List Shows", show)
	return nil
}
