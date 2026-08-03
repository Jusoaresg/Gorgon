package api

import (
	"log/slog"
	"strconv"

	"github.com/jusoaresg/gorgon/pkg/schemas"

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
func (h *Handler) GetShowAliases(c echo.Context) error {
	h.Logger.Info("Received request to Get Show Aliases", slog.String("endpoint", "/database/show/aliases"), slog.String("method", "get"))

	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		h.Logger.Error("failed to convert id to int", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error while converting id to int")
		return err
	}

	aliases, err := h.ShowAliasRepo.ListByShowID(idInt64)
	if err != nil {
		h.Logger.Error("Error while fetching show aliases from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show aliases")
		return err
	}

	h.Logger.Info("Successfully fetched show aliases")
	schemas.SendSuccess(c, "Get Show Aliases", aliases)
	return nil
}
