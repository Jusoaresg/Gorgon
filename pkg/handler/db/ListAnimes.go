package handler

import (
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary List Animes
// @Description List all added animes
// @Tags Database/Anime
// @Produce json
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/anime [get]
func ListAnimes(c *gin.Context) {
	animes := []schemas.Anime{}

	baseService := services.NewBaseService()

	if err := baseService.List(&animes); err != nil {
		return
	}

	schemas.SendSucess(c, "ListAnimes", animes)
}
