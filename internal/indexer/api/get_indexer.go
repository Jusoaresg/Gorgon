package api

import (
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Get Indexer
// @Description Get indexer
// @Tags Database/Indexer
// @Produce json
// @Param id path string true "Indexer identification"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/indexer/{id} [get]
func (h *Handler) GetIndexer(c echo.Context) error {
	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 400, "Failed to bind request body")
		return err
	}

	indexer, err := h.IndexerRepo.GetById(int(request.Id))
	if err != nil {
		return err
	}

	schemas.SendSuccess(c, "Get Indexer", indexer)
	return nil
}
