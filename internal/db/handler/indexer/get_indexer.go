package indexer

import (
	"gorgon/internal/db/repository"
	"gorgon/pkg/schemas"

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
func GetIndexer(c echo.Context) error {
	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 400, "Failed to bind request body")
		return err
	}

	indexerRepo := repository.NewIndexerRepository()

	indexer, err := indexerRepo.GetById(int(request.Id))
	if err != nil {
		return err
	}

	schemas.SendSucess(c, "List Indexers", indexer)
	return nil
}
