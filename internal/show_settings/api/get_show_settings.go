package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"strconv"

	"github.com/jusoaresg/gorgon/internal/show_settings/schema"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Get Show Settings
// @Description Get Show Settings
// @Tags Database/ShowSettings
// @Produce json
// @Param id path int true "Show ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show-settings/{id} [get]
func (h *Handler) GetShowSettings(c echo.Context) error {
	h.Logger.Info("Received request to Get Show Settings", slog.String("endpoint", "/database/show-settings"), slog.String("method", "get"))

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Logger.Error("Failed to convert id to int64", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error while converting id to int64")
		return err
	}

	settings, err := h.ShowSettingsRepo.GetByShowID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			schemas.SendSuccess(c, "Get Show Settings", schema.ShowSettingsDto{
				FilterProfileID: nil,
				UseAliases:      true,
				OnlyLatin:       true,
			})
			return nil
		}
		h.Logger.Error("Error while fetching show settings", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show settings")
		return err
	}

	h.Logger.Info("Successfully fetched show settings", slog.Int64("show_id", id))
	schemas.SendSuccess(c, "Get Show Settings", schema.ToShowSettingsDto(settings))
	return nil
}
