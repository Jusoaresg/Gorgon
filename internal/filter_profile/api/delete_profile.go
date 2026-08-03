package api

import (
	"errors"
	"log/slog"
	"strconv"

	filterProfileRepository "github.com/jusoaresg/gorgon/internal/filter_profile/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Delete Filter Profile
// @Description Delete Filter Profile
// @Tags Database/FilterProfile
// @Produce json
// @Param id path int true "Filter Profile ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/filter-profile/{id} [delete]
func (h *Handler) DeleteFilterProfile(c echo.Context) error {
	h.Logger.Info("Received request to Delete Filter Profile", slog.String("endpoint", "/database/filter-profile"), slog.String("method", "delete"))

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Logger.Error("Failed to convert id to int64", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error while converting id to int64")
		return err
	}

	err = h.FilterProfileRepo.Delete(id)
	if err != nil {
		if errors.Is(err, filterProfileRepository.ErrFilterProfileNotFound) {
			schemas.SendError(c, 404, "Filter profile not found")
			return err
		}
		h.Logger.Error("Error while deleting filter profile", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while deleting filter profile")
		return err
	}

	h.Logger.Info("Successfully deleted filter profile", slog.Int64("id", id))
	schemas.SendSuccess(c, "Delete Filter Profile", map[string]int64{"id": id})
	return nil
}
