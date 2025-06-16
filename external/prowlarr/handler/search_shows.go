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

// @Summary Search For Shows
// @Description Search shows on prowlarr indexers
// @Tags Prowlarr/Search
// @Produce json
// @Param request body schema.SearchRequest true "Request Body"
// @Success 200 {object} []schema.SearchResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /prowlarr/search [post]
func SearchShows(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Search Shows", slog.String("endpoint", "/api/v1/prowlarr/search"), slog.String("method", "POST"))

	var request schema.SearchRequest
	if err := c.Bind(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to bind request body")
		return err
	}

	searchService := service.NewProwlarrSearchService(logger)

	var response []schema.SearchResponse
	if err := searchService.Search(&request, &response); err != nil {
		logger.Error("Error while search all shows from prowlarr", slog.String("error", err.Error()))
		schemas.SendError(c, 500, fmt.Sprintf("Error while searching all animes prowlarr: %s", err.Error()))
		return err
	}

	logger.Info("Search Shows request successfully", slog.Any("response", response))
	schemas.SendSucess(c, "Search Shows", response)
	return nil
}
