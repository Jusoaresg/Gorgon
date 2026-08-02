package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"strconv"

	"github.com/jusoaresg/gorgon/internal/filter_profile/schema"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Get Filter Profile
// @Description Get Filter Profile
// @Tags Database/FilterProfile
// @Produce json
// @Param id path int true "Filter Profile ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/filter-profile/{id} [get]
func (h *Handler) GetFilterProfile(c echo.Context) error {
	h.Logger.Info("Received request to Get Filter Profile", slog.String("endpoint", "/database/filter-profile"), slog.String("method", "get"))

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Logger.Error("Failed to convert id to int64", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error while converting id to int64")
		return err
	}

	profile, patterns, err := h.FilterProfileRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			schemas.SendError(c, 404, "Filter profile not found")
			return err
		}
		h.Logger.Error("Error while fetching filter profile", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching filter profile")
		return err
	}

	h.Logger.Info("Successfully fetched filter profile", slog.Int64("id", id))
	schemas.SendSuccess(c, "Get Filter Profile", schema.ToProfileDto(profile, patterns))
	return nil
}
