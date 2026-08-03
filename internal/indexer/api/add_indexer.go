package api

import (
	"encoding/json"

	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/internal/indexer/model"
	"github.com/jusoaresg/gorgon/pkg/schemas"

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
func (h *Handler) AddIndexer(c echo.Context) error {
	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	var response schema.IndexerResponse
	if err := h.ProwlarrIndexerSvc.GetIndexer(int(request.Id), &response); err != nil {
		schemas.SendError(c, 500, "Failed to get indexer")
		return err
	}

	indexerUrlsJSON, err := json.Marshal(response.IndexerUrls)
	if err != nil {
		schemas.SendError(c, 500, "Failed to marshal indexer URLs")
		return err
	}

	indexer := model.Indexer{
		IndexerID:      response.Id,
		Name:           response.Name,
		DefinitionName: response.DefinitionName,
		IndexerUrls:    string(indexerUrlsJSON),
		Language:       response.Language,
	}

	if err := h.IndexerRepo.Create(indexer); err != nil {
		return err
	}

	schemas.SendSuccess(c, "Add Indexer", indexer)
	return nil
}
