package api

import (
	"fmt"

	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Delete Indexer
// @Description Delete Indexer from db
// @Tags Database/Indexer
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/indexer [delete]
func (h *Handler) DeleteIndexer(c echo.Context) error {
	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	if err := h.IndexerRepo.DeleteById(int(request.Id)); err != nil {
		schemas.SendError(c, 500, fmt.Sprintf("Error while deleting indexer: %s", err.Error()))
		return err
	}

	schemas.SendSuccess(c, "Delete Indexer", "")
	return nil
}
