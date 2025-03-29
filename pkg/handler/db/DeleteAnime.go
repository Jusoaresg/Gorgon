package handler

import (
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Delete Anime
// @Description Delete Anime from List
// @Tags Database/Anime
// @Produce json
// @Param request body schemas.DeleteAnimeRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/anime [delete]
func DeleteAnime(c *gin.Context) {
	request := schemas.DeleteAnimeRequest{}
	c.BindJSON(&request)

	id := request.Id
	if id == "" {
		return
	}

	anime := schemas.Anime{}

	baseService := services.NewBaseService()

	if err := baseService.Get(&anime, id).Error; err != nil {
		return
	}

	if err := baseService.Delete(string(anime.ID), schemas.Anime{}).Error; err != nil {
		return
	}

	schemas.SendSucess(c, "DeleteAnime", anime)
}
