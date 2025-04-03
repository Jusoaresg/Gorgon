package indexer

import (
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
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
func ListIndexers(c *gin.Context) {
	indexers := []model.Indexer{}

	baseService := services.NewBaseService()

	if err := baseService.List(&indexers); err != nil {
		return
	}

	schemas.SendSucess(c, "List Indexers", indexers)
}
