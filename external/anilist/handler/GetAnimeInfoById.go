package handler

import (
	"gorgon/config"
	"gorgon/external/anilist/service"
	"gorgon/pkg/schemas"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Get Anime info By id
// @Description Search anime info by id
// @Tags Anilist
// @Accept json
// @Produce json
// @Param request body schema.AnimeGetInfoByIdRequest true "Request Body"
// @Success 200 {object} dtos.AnimeDto
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /anilist/get-info [post]
func GetAnimeInfoById(c *gin.Context) {
	logger := config.GetLogger()
	logger.Info("Received requets to Get Anime Info By ID", slog.String("endpoint", "/api/v1/anilist/get-info"), slog.String("method", "POST"))

	var request schemas.IdRequest
	if err := c.BindJSON(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("endpoint", "/api/v1/anilist/get-info"), slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Invalid request format")
		return
	}

	anilistService := service.NewAnimeListService(logger)

	response, err := anilistService.GetAnimeInfoById(request.Id)
	if err != nil {
		logger.Error("Failed to Get Anime info By Id", slog.String("endpoint", "/api/v1/anilist/get-info"), slog.String("method", "POST"), slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to get anime info")
		return
	}

	logger.Info("Successfully find anime info by id", slog.Int("id", request.Id), slog.String("title", response.Title.Romaji))
	schemas.SendSucess(c, "GetAnimeInfoById", response)
}
