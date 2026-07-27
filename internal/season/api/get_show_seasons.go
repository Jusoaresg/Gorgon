package api

import (
	"log/slog"
	"strconv"

	"github.com/jusoaresg/gorgon/pkg/schemas"

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
func (h *Handler) GetShowSeasons(c echo.Context) error {
	h.Logger.Info("Received request to Get Show Seasons", slog.String("endpoint", "/database/show/seasons/:id"), slog.String("method", "get"))

	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		schemas.SendError(c, 400, "Error while converting id to int")
		return err
	}

	seasons, err := h.SeasonRepo.ListByShowId(idInt64)
	if err != nil {
		h.Logger.Error("Error while fetching show seasons from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show seasons")
		return err
	}

	h.Logger.Info("Successfully fetched show seasons")
	schemas.SendSuccess(c, "Get Show Seasons", seasons)
	return nil
}
