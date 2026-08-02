package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/internal/filter_profile/schema"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary List Filter Profiles
// @Description List Filter Profiles
// @Tags Database/FilterProfile
// @Produce json
// @Success 200 {object} schemas.DefaultResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/filter-profile [get]
func (h *Handler) ListFilterProfiles(c echo.Context) error {
	h.Logger.Info("Received request to List Filter Profiles", slog.String("endpoint", "/database/filter-profile"), slog.String("method", "get"))

	profiles, err := h.FilterProfileRepo.List()
	if err != nil {
		h.Logger.Error("Error while listing filter profiles", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while listing filter profiles")
		return err
	}

	h.Logger.Info("Successfully listed filter profiles", slog.Int("count", len(profiles)))
	schemas.SendSuccess(c, "List Filter Profiles", schema.ToProfileListDto(profiles))
	return nil
}
