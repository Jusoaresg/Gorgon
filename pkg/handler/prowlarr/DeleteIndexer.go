package prowlarr

import (
	"fmt"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Delete Indexer
// @Description Delete Indexer from db
// @Tags Prowlarr
// @Produce json
// @Param request body schemas.RemoveIndexerRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /prowlarr/indexers [delete]
func DeleteIndexer(c *gin.Context) {
	request := schemas.RemoveIndexerRequest{}
	c.BindJSON(&request)

	baseService := services.NewBaseService()

	if err := baseService.DeletePermanently(request.Id, &schemas.Indexer{}); err != nil {
		schemas.SendError(c, 500, fmt.Sprintf("Error while deleting indexer: %w", err))
	}
}
