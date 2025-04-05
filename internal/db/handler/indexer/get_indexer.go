package indexer

import (
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

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
	// id := c.Param("id")
	// idInt, err := strconv.Atoi(id)
	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 400, "Failed to bind request body")
		return err
	}

	indexer := model.Indexer{}

	baseService := services.NewBaseService()

	if err := baseService.Get(&indexer, request.Id); err != nil {
		return err
	}

	schemas.SendSucess(c, "List Indexers", indexer)
	return nil
}
