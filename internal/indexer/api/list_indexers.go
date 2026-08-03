package api

import (
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary List Indexers
// @Description List all added indexers
// @Tags Database/Indexer
// @Produce json
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/indexer [get]
func (h *Handler) ListIndexers(c echo.Context) error {
	indexers, err := h.IndexerRepo.List()
	if err != nil {
		return err
	}

	schemas.SendSuccess(c, "List Indexers", indexers)
	return nil
}
