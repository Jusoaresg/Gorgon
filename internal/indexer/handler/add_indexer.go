package indexer

import (
	"encoding/json"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/external/prowlarr/service"
	"github.com/jusoaresg/gorgon/internal/indexer/model"
	"github.com/jusoaresg/gorgon/internal/indexer/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
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
func AddIndexer(c echo.Context) error {
	logger := config.GetLogger()

	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	indexerRepo := repository.NewIndexerRepository()
	prowlarrIndexerService := service.NewProwlarrIndexerService(logger)

	var response schema.IndexerResponse
	if err := prowlarrIndexerService.GetIndexer(int(request.Id), &response); err != nil {
		c.JSON(500, gin.H{"error": "Failed to get indexer"})
		return err
	}

	indexerUrlsJSON, err := json.Marshal(response.IndexerUrls)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to marshal indexer URLs"})
		return err
	}

	indexer := model.Indexer{
		IndexerID:      response.Id,
		Name:           response.Name,
		DefinitionName: response.DefinitionName,
		IndexerUrls:    string(indexerUrlsJSON),
		Language:       response.Language,
	}

	if err := indexerRepo.Create(indexer); err != nil {
		return err
	}

	schemas.SendSuccess(c, "Add Indexer", indexer)
	return nil
}
