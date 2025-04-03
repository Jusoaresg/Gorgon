package handler

import (
	"gorgon/config"
	"gorgon/external/anilist/schema"
	"gorgon/external/anilist/service"
	"gorgon/pkg/schemas"
	"log/slog"

	"github.com/gin-gonic/gin"

	"net/http"
)

// @BasePath /api/v1

// @Summary Find Anime by name
// @Description Search anime by name
// @Tags Anilist
// @Accept json
// @Produce json
// @Param request body schema.AnimeTitleIdByNameRequest true "Request Body"
// @Success 200 {object} schema.AnimeTitleIdByNameResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /anilist [post]
func FindAnimeIdAnilist(c *gin.Context) {
	logger := config.GetLogger()
	logger.Info("Received request to FindAnimeIdAnilist(Find Anime by name)", slog.String("endpoint", "/api/v1/anilist"))

	request := schema.AnimeTitleIdByNameRequest{}
	if err := c.BindJSON(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("endpoint", "/api/v1/anilist"), slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Invalid request body")
	}

	if request.Name == "" {
		logger.Warn("Anime Name is required but not provided", slog.String("request_name", request.Name))
		schemas.SendError(c, 500, "Anime Name is required")
		return
	}

	anilistService := service.NewAnimeListService(logger)

	response, err := anilistService.FindAnimeIdAnilist(&request)
	if err != nil {
		logger.Error("Anilist service returned an error", slog.String("error", err.Error()))
		return
	}

	logger.Info("Successfully fetched anime data from anilist", slog.String("anime_name", request.Name), slog.Int("media_count", len(response.Data.Page.Media)))
	c.JSON(http.StatusOK, gin.H{
		"media": response.Data.Page.Media,
	})
}
