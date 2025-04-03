package anime

import (
	"gorgon/config"
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Delete Anime
// @Description Delete Anime from List
// @Tags Database/Anime
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/anime [delete]
func DeleteAnime(c *gin.Context) {
	logger := config.GetLogger()
	logger.Info("Received request to Delete Anime", slog.String("endpoint", "/api/v1/database/anime"))

	var request schemas.IdRequest
	if err := c.BindJSON(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to bind request body")
		return
	}

	id := request.Id
	if request.Id == 0 {
		logger.Error("Invalid ID", slog.String("error", "ID cannot be 0"))
		schemas.SendError(c, 400, "ID cannot be 0")
		return
	}

	anime := model.Anime{}

	baseService := services.NewBaseService()

	if err := baseService.Get(&anime, id); err != nil {
		logger.Error("Error while retrieving anime data from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while retrieving anime data from database")
		return
	}

	if err := baseService.Delete(anime.Aid, model.Anime{}); err != nil {
		logger.Error("Error while deleting anime from databases", slog.Int("aid", request.Id), slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while deleting anime from database")
		return
	}

	schemas.SendSucess(c, "DeleteAnime", anime)
}
