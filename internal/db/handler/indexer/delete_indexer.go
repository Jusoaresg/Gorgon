package indexer

import (
	"fmt"
	"gorgon/internal/db/repository"
	"gorgon/pkg/schemas"

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
func DeleteIndexer(c echo.Context) error {
	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	indexerRepo := repository.NewIndexerRepository()

	if err := indexerRepo.DeleteById(int(request.Id)); err != nil {
		schemas.SendError(c, 500, fmt.Sprintf("Error while deleting indexer: %s", err.Error()))
		return err
	}

	schemas.SendSucess(c, "Delete Indexer", "")
	return nil
}
