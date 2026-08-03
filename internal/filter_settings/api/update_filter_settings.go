package api

import (
	"log/slog"
	"strconv"

	"github.com/jusoaresg/gorgon/internal/filter_settings/schema"
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

	var request schema.FilterSettingsDto
	if err := c.Bind(&request); err != nil {
		h.Logger.Error("Failed to bind request", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error binding request")
		return err
	}

	settings, err := h.FilterSettingsRepo.Get()
	if err != nil {
		h.Logger.Error("Error while fetching filter settings", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching filter settings")
		return err
	}

	if request.DefaultFilterProfileID != nil {
		if *request.DefaultFilterProfileID == "" {
			settings.DefaultFilterProfileID = nil
		} else {
			id, err := strconv.ParseInt(*request.DefaultFilterProfileID, 10, 64)
			if err != nil {
				h.Logger.Error("Invalid default filter profile id", slog.String("error", err.Error()))
				schemas.SendError(c, 400, "Invalid default filter profile id")
				return err
			}
			settings.DefaultFilterProfileID = &id
		}
	}

	if request.UseAliases != nil {
		settings.UseAliases = *request.UseAliases
	}

	if request.OnlyLatin != nil {
		settings.OnlyLatin = *request.OnlyLatin
	}

	if err := h.FilterSettingsRepo.Save(settings); err != nil {
		h.Logger.Error("Error while updating filter settings", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while updating filter settings")
		return err
	}

	h.Logger.Info("Successfully updated filter settings")
	schemas.SendSuccess(c, "Update Filter Settings", settings)
	return nil
}
