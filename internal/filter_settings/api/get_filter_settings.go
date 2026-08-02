package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Get Filter Settings
// @Description Get Global Filter Settings
// @Tags Database/FilterSettings
// @Produce json
// @Success 200 {object} schemas.DefaultResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/filter-settings [get]
func (h *Handler) GetFilterSettings(c echo.Context) error {
	h.Logger.Info("Received request to Get Filter Settings", slog.String("endpoint", "/database/filter-settings"), slog.String("method", "get"))

	settings, err := h.FilterSettingsRepo.Get()
	if err != nil {
		h.Logger.Error("Error while fetching filter settings", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching filter settings")
		return err
	}

	h.Logger.Info("Successfully fetched filter settings")
	schemas.SendSuccess(c, "Get Filter Settings", settings)
	return nil
}
