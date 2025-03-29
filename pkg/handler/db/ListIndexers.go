package handler

import (
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary List Added Indexers
// @Description List all added indexers
// @Tags Database/Indexer
// @Produce json
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/indexer [get]
func ListIndexers(c *gin.Context) {
	indexers := []schemas.Indexer{}

	baseService := services.NewBaseService()

	if err := baseService.List(&indexers); err != nil {
		return
	}

	schemas.SendSucess(c, "List Indexers", indexers)
}

// @BasePath /api/v1

// @Summary Get Indexers
// @Description Get indexers
// @Tags Database/Indexer
// @Produce json
// @Param id path string true "Indexer identification"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/indexer/{id} [get]
func GetIndexer(c *gin.Context) {
	id := c.Param("id")

	indexer := schemas.Indexer{}

	baseService := services.NewBaseService()

	if err := baseService.Get(&indexer, id); err != nil {
		return
	}

	schemas.SendSucess(c, "List Indexers", indexer)
}
