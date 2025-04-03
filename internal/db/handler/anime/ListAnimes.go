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

// @Summary List Animes
// @Description List all animes
// @Tags Database/Anime
// @Produce json
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/anime [get]
func ListAnimes(c *gin.Context) {
	logger := config.GetLogger()
	logger.Info("Received request to list animes", slog.String("endpoint", "/database/anime"), slog.String("method", "get"))

	animes := []model.Anime{}

	baseService := services.NewBaseService()

	if err := baseService.List(&animes); err != nil {
		logger.Error("Error while fetching animes from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching animes")
		return
	}

	logger.Info("Successfully fetched animes", slog.Int("count", len(animes)))
	schemas.SendSucess(c, "ListAnimes", animes)
}
