package handler

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/prowlarr/schema"
	"gorgon/external/prowlarr/service"
	"gorgon/pkg/schemas"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Get All Indexer
// @Description Get All Prowlarr indexers
// @Tags Prowlarr/Indexer
// @Accept json
// @Produce json
// @Success 200 {object} []schema.IndexerResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /prowlarr/indexer [get]
func GetIndexer(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Get Indexer", slog.String("endpoint", "/api/v1/prowlarr/indexer"))

	prowlarrIndexerService := service.NewProwlarrIndexerService(logger)

	var indexers []schema.IndexerResponse
	if err := prowlarrIndexerService.GetIndexers(&indexers); err != nil {
		logger.Error("Error getting list of indexers", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error getting list of indexers: %s", err.Error()))
		return err
	}

	logger.Info("Get Indexer request successfully", slog.Any("response", indexers))
	schemas.SendSucess(c, "Get Indexer", indexers)
	return nil
}
