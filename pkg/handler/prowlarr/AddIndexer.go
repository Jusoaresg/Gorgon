package prowlarr

import (
	"encoding/json"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// @BasePath /api/v1

// @Summary Add Indexer
// @Description Add prowlarr indexer
// @Tags Prowlarr
// @Accept json
// @Produce json
// @Param request body schemas.AddIndexerRequest true "Request Body"
// @Success 200 {object} []schemas.AddIndexerRequest
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /prowlarr/indexers [post]
func AddIndexer(c *gin.Context) {
	request := schemas.AddIndexerRequest{}
	c.BindJSON(&request)

	baseService := services.NewBaseService()

	prowlarrIndexerService := services.NewProwlarrIndexerService()

	var response schemas.IndexerResponse
	if err := prowlarrIndexerService.GetIndexer(request.IndexerId, &response); err != nil {
		c.JSON(500, gin.H{"error": "Failed to get indexer"})
		return
	}

	indexerUrlsJSON, err := json.Marshal(response.IndexerUrls)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal indexer URLs"})
		return
	}

	indexer := schemas.Indexer{
		IndexerId:      response.Id,
		Name:           response.Name,
		DefinitionName: response.DefinitionName,
		IndexerUrls:    datatypes.JSON(indexerUrlsJSON),
		Language:       response.Language,
	}

	baseService.Add(&indexer)

	schemas.SendSucess(c, "Add Indexer", indexer)
}
