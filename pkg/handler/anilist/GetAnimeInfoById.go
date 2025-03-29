package anilist

import (
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Get Anime info By id
// @Description Search anime info by id
// @Tags Anilist
// @Accept json
// @Produce json
// @Param request body schemas.AnimeGetInfoByIdRequest true "Request Body"
// @Success 200 {object} schemas.AnimeGetInfoByIdResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /anilist/get-info [post]
func GetAnimeInfoById(c *gin.Context) {
	request := schemas.AnimeGetInfoByIdRequest{}
	c.BindJSON(&request)

	anilistService := services.NewAnimeListService()

	response := anilistService.GetAnimeInfoById(&request)

	c.JSON(http.StatusOK, gin.H{
		"media": response.Data.Page.Media,
	})
}
