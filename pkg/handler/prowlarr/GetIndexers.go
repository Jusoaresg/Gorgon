package prowlarr

import (
	"fmt"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Get Indexers
// @Description Get Prowlarr indexers
// @Tags Prowlarr
// @Accept json
// @Produce json
// @Success 200 {object} []schemas.IndexerResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /prowlarr/indexers [get]
func GetIndexers(c *gin.Context) {

	prowlarrIndexerService := services.NewProwlarrIndexerService()

	var indexers []schemas.IndexerResponse
	if err := prowlarrIndexerService.GetIndexers(&indexers); err != nil {
		schemas.SendError(c, 500, fmt.Sprintf("Error getting indexers: %w", err))
		return
	}

	schemas.SendSucess(c, "Get Indexers", indexers)
}
