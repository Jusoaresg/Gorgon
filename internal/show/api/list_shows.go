package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/pkg/schemas"

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
func (h *Handler) ListShows(c echo.Context) error {
	h.Logger.Info("Received request to List Shows", slog.String("endpoint", "/database/show"), slog.String("method", "get"))

	search := c.QueryParam("search")
	status := c.QueryParam("status")
	sort := c.QueryParam("sort")
	_ = sort

	show, err := h.ShowRepo.ListFiltered(search, status)
	if err != nil {
		h.Logger.Error("Error while fetching shows from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching shows")
		return err
	}

	h.Logger.Info("Successfully fetched shows", slog.Int("count", len(show)))
	schemas.SendSuccess(c, "List Shows", show)
	return nil
}
