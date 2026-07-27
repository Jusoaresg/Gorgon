package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/pkg/schemas"

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
func (h *Handler) ListFullShows(c echo.Context) error {
	h.Logger.Info("Received request to List Full Shows", slog.String("endpoint", "/database/show/full"), slog.String("method", "get"))

	shows, err := h.AggregatorService.ListFullShows()
	if err != nil {
		h.Logger.Error("Error while fetching shows from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching shows")
		return err
	}

	h.Logger.Info("Successfully fetched full shows", slog.Int("count", len(shows)))
	schemas.SendSuccess(c, "List Full Shows", shows)
	return nil
}
