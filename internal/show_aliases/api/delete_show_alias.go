package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"strconv"

	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Delete Show Alias
// @Description Delete a custom alias for a show
// @Tags Database/Aliases
// @Produce json
// @Param id path int true "Show ID"
// @Param aliasId path int true "Alias ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 404 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/{id}/alias/{aliasId} [delete]
func (h *Handler) DeleteShowAlias(c echo.Context) error {
	h.Logger.Info("Received request to Delete Show Alias", slog.String("endpoint", "/database/show/alias"), slog.String("method", "delete"))

	showID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Logger.Error("Failed to convert id to int64", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error while converting id to int64")
		return err
	}

	aliasID, err := strconv.ParseInt(c.Param("aliasId"), 10, 64)
	if err != nil {
		h.Logger.Error("Failed to convert aliasId to int64", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error while converting aliasId to int64")
		return err
	}

	alias, err := h.ShowAliasRepo.GetByID(aliasID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			schemas.SendError(c, 404, "Alias not found")
			return err
		}
		h.Logger.Error("Error while fetching show alias", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show alias")
		return err
	}

	if alias.ShowID != showID {
		schemas.SendError(c, 400, "Alias does not belong to the given show")
		return nil
	}

	if alias.Source != "user" {
		schemas.SendError(c, 400, "Only user-created aliases can be deleted")
		return nil
	}

	if err := h.ShowAliasRepo.DeleteByID(aliasID); err != nil {
		h.Logger.Error("Error while deleting show alias", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while deleting show alias")
		return err
	}

	h.Logger.Info("Successfully deleted show alias", slog.Int64("alias_id", aliasID))
	schemas.SendSuccess(c, "Delete Show Alias", map[string]int64{"id": aliasID})
	return nil
}
