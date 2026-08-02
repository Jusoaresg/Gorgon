package api

import (
	"log/slog"

	filterSettingsModel "github.com/jusoaresg/gorgon/internal/filter_settings/model"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Update Filter Settings
// @Description Update Global Filter Settings
// @Tags Database/FilterSettings
// @Accept json
// @Produce json
// @Param request body filterSettingsModel.FilterSettings true "Filter settings data"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/filter-settings [patch]
func (h *Handler) UpdateFilterSettings(c echo.Context) error {
	h.Logger.Info("Received request to Update Filter Settings", slog.String("endpoint", "/database/filter-settings"), slog.String("method", "patch"))

	var request filterSettingsModel.FilterSettings
	if err := c.Bind(&request); err != nil {
		h.Logger.Error("Failed to bind request", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error binding request")
		return err
	}

	if err := h.FilterSettingsRepo.Save(request); err != nil {
		h.Logger.Error("Error while updating filter settings", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while updating filter settings")
		return err
	}

	h.Logger.Info("Successfully updated filter settings")
	schemas.SendSuccess(c, "Update Filter Settings", request)
	return nil
}
