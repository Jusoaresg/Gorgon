package api

import (
	"log/slog"
	"strconv"
	"strings"

	showAliasModel "github.com/jusoaresg/gorgon/internal/show_aliases/model"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type addShowAliasRequest struct {
	Alias string `json:"alias"`
}

// @BasePath /api/v1

// @Summary Add Show Alias
// @Description Add a custom alias for a show
// @Tags Database/Aliases
// @Accept json
// @Produce json
// @Param id path int true "Show ID"
// @Param request body addShowAliasRequest true "Alias data"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/{id}/alias [post]
func (h *Handler) AddShowAlias(c echo.Context) error {
	h.Logger.Info("Received request to Add Show Alias", slog.String("endpoint", "/database/show/alias"), slog.String("method", "post"))

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.Logger.Error("Failed to convert id to int64", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Error while converting id to int64")
		return err
	}

	var request addShowAliasRequest
	if err := c.Bind(&request); err != nil {
		h.Logger.Error("Failed to bind request", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error binding request")
		return err
	}

	if strings.TrimSpace(request.Alias) == "" {
		schemas.SendError(c, 400, "Alias is required")
		return nil
	}

	aliasID, err := h.ShowAliasRepo.Create(showAliasModel.ShowAlias{
		ShowID: id,
		Alias:  strings.TrimSpace(request.Alias),
		Source: "user",
	})
	if err != nil {
		h.Logger.Error("Error while adding show alias", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while adding show alias")
		return err
	}

	h.Logger.Info("Successfully added show alias", slog.Int64("show_id", id), slog.Int64("alias_id", aliasID))
	schemas.SendSuccess(c, "Add Show Alias", map[string]int64{"id": aliasID})
	return nil
}
