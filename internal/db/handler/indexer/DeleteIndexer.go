package indexer

import (
	"fmt"
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
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
func DeleteIndexer(c *gin.Context) {
	var request schemas.IdRequest
	c.BindJSON(&request)

	baseService := services.NewBaseService()

	if err := baseService.DeletePermanently(request.Id, &model.Indexer{}); err != nil {
		schemas.SendError(c, 500, fmt.Sprintf("Error while deleting indexer: %s", err.Error()))
		return
	}
}
