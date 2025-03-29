package anilist

import (
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"

	"github.com/gin-gonic/gin"

	"net/http"
)

// @BasePath /api/v1

// @Summary Find Anime by name
// @Description Search anime by name
// @Tags Anilist
// @Accept json
// @Produce json
// @Param request body schemas.AnimeTitleIdByNameRequest true "Request Body"
// @Success 200 {object} schemas.AnimeTitleIdByNameResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /anilist [post]
func FindAnimeIdAnilist(c *gin.Context) {

	request := schemas.AnimeTitleIdByNameRequest{}
	c.BindJSON(&request)

	if request.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anime name is required"})
		return
	}

	anilistService := services.NewAnimeListService()

	response := anilistService.FindAnimeIdAnilist(&request)

	c.JSON(http.StatusOK, gin.H{
		"media": response.Data.Page.Media,
	})
}
