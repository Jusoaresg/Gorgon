package episode

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/show_aliases/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary List Show Aliases
// @Description List Show Aliases
// @Tags Database/Aliases
// @Produce json
// @Param id path int true "Show ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/aliases/{id} [get]
func GetShowEpisodes(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Get Show Aliases", slog.String("endpoint", "/database/show/aliases"), slog.String("method", "get"))

	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Error("failed to convert id to int", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error while converting id to int")
		return err
	}

	showAliasRepo := repository.NewShowAliasesRepository(config.GetSQLite())
	aliases, err := showAliasRepo.ListByShowID(idInt64)
	if err != nil {
		logger.Error("Error while fetching show aliases from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show aliases")
		return err
	}

	logger.Info("Successfully fetched show aliases")
	schemas.SendSuccess(c, "Get Show Aliases", aliases)
	return nil
}
