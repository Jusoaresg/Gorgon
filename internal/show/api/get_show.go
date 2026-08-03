package api

import (
	"log/slog"
	"strconv"

	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Get Show
// @Description Getshow
// @Tags Database/Show
// @Produce json
// @Param id path int true "Show ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/{id} [get]
func (h *Handler) GetShow(c echo.Context) error {
	h.Logger.Info("Received request to Get Show", slog.String("endpoint", "/database/show/:id"), slog.String("method", "get"))

	id := c.Param("id")
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		schemas.SendError(c, 400, "Error parsing id to int64")
		return err
	}

	//TODO: Maybe returning the episodes and seasons ?
	show, err := h.ShowRepo.GetById(id64)
	if err != nil {
		h.Logger.Error("Error while fetching show from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show")
		return err
	}

	h.Logger.Info("Successfully fetched show", slog.Any("Show", show))
	schemas.SendSuccess(c, "Get Show", show)
	return nil
}
