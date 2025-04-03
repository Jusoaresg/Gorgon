package indexer

import (
	"encoding/json"
	"gorgon/config"
	"gorgon/external/prowlarr/schema"
	"gorgon/external/prowlarr/service"
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

// @BasePath /api/v1

// @Summary Add Indexer
// @Description Add Indexer to db
// @Tags Database/Indexer
// @Accept json
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/indexer [post]
func AddIndexer(c *gin.Context) {
	logger := config.GetLogger()

	var request schemas.IdRequest
	c.BindJSON(&request)

	baseService := services.NewBaseService()

	prowlarrIndexerService := service.NewProwlarrIndexerService(logger)

	var response schema.IndexerResponse
	if err := prowlarrIndexerService.GetIndexer(request.Id, &response); err != nil {
		c.JSON(500, gin.H{"error": "Failed to get indexer"})
		return
	}

	indexerUrlsJSON, err := json.Marshal(response.IndexerUrls)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal indexer URLs"})
		return
	}

	indexer := model.Indexer{
		IndexerId:      response.Id,
		Name:           response.Name,
		DefinitionName: response.DefinitionName,
		IndexerUrls:    datatypes.JSON(indexerUrlsJSON),
		Language:       response.Language,
	}

	baseService.Add(&indexer)

	schemas.SendSucess(c, "Add Indexer", indexer)
}
