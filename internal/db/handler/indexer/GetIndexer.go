package indexer

import (
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
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
func GetIndexer(c *gin.Context) {
	// id := c.Param("id")
	// idInt, err := strconv.Atoi(id)
	var request schemas.IdRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	indexer := model.Indexer{}

	baseService := services.NewBaseService()

	if err := baseService.Get(&indexer, request.Id); err != nil {
		return
	}

	schemas.SendSucess(c, "List Indexers", indexer)
}
