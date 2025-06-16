package season

import (
	"gorgon/config"
	"gorgon/internal/db/repository"
	"gorgon/pkg/schemas"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary List Show Seasons
// @Description List Show Seasons
// @Tags Database/Seasons
// @Produce json
// @Param id path int true "Show ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/seasons/{id} [get]
func GetShowSeasons(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Get Show Seasons", slog.String("endpoint", "/database/show/seasons/:id"), slog.String("method", "get"))

	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		schemas.SendError(c, 400, "Error while converting id to int")
		return err
	}

	seasonRepo := repository.NewSeasonRepository()
	seasons, err := seasonRepo.ListByShowId(idInt64)
	if err != nil {
		logger.Error("Error while fetching show seasons from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show seasons")
		return err
	}

	logger.Info("Successfully fetched show seasons")
	schemas.SendSucess(c, "Get Show Seasons", seasons)
	return nil
}
