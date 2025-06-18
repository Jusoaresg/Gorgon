package episode

import (
	"gorgon/config"
	"gorgon/internal/db/repository"
	"gorgon/pkg/schemas"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary List Show Episodes
// @Description List Show Episodes
// @Tags Database/Episodes
// @Produce json
// @Param id path int true "Show ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/episode/{id} [get]
func GetShowEpisodes(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Get Show Episodes", slog.String("endpoint", "/database/show/episode/:id"), slog.String("method", "get"))

	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		schemas.SendError(c, 400, "Error while converting id to int")
		return err
	}

	episodeRepo := repository.NewEpisodeRepository(config.GetSQLite())
	episodes, err := episodeRepo.ListByShowID(idInt64)
	if err != nil {
		logger.Error("Error while fetching show episodes from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show episodes")
		return err
	}

	logger.Info("Successfully fetched show episodes")
	schemas.SendSucess(c, "Get Show Episodes", episodes)
	return nil
}
