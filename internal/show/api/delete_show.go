package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jusoaresg/gorgon/internal/show/model"
	"github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Delete Show
// @Description Delete Show from list
// @Tags Database/Show
// @Produce json
// @Param id path int true "Show Id"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/{id} [delete]
func (h *Handler) DeleteShow(c echo.Context) error {
	h.Logger.Info("Received request to Delete Show", slog.String("endpoint", "/api/v1/database/show"))

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.Logger.Error("Failed parse id", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Failed to bind request body")
		return err
	}

	if id <= 0 {
		h.Logger.Error("Invalid ID", slog.String("error", "ID must be greater than 0"))
		schemas.SendError(c, 400, "ID must be greater than 0")
		return errors.New("ID must be greater than 0")
	}

	show := model.Show{}

	err = h.ShowRepo.DeleteById(id)
	if err != nil {
		if errors.Is(err, repository.ErrShowNotFound) {
			schemas.SendError(c, 404, "Show not found")
			return err
		}
		schemas.SendError(c, 500, "Internal server error")
		return err
	}

	isHtmx := c.Request().Header.Get("HX-Request") == "true"
	if isHtmx {
		c.Response().Header().Set("HX-Redirect", "/")
		return c.NoContent(http.StatusOK)
	}

	schemas.SendSuccess(c, "DeleteShow", show)
	return nil
}
