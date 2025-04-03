package indexer

import (
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

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
func ListIndexers(c echo.Context) error {
	indexers := []model.Indexer{}

	baseService := services.NewBaseService()

	if err := baseService.List(&indexers); err != nil {
		return err
	}

	schemas.SendSucess(c, "List Indexers", indexers)
	return nil
}
