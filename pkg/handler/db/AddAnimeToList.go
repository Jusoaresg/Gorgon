package handler

import (
	"encoding/json"
	"fmt"
	"gorgon/config"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Add Anime
// @Description Add anime to list
// @Tags Database/Anime
// @Accept json
// @Produce json
// @Param request body schemas.AddAnimeToListRequest true "Request Body"
// @Success 200 {object} schemas.AddAnimeToListResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/anime [post]
func AddAnimeToList(c *gin.Context) {
	request := schemas.AddAnimeToListRequest{}
	c.BindJSON(&request)

	// reqBodyJSON, err := json.Marshal(request)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode request body"})
	// 	return
	// }

	apiService := services.NewAPIService(fmt.Sprintf("http://127.0.0.1:%s", config.Port))

	var response schemas.AddAnimeToListResponse
	apiService.Post("/api/v1/anilist/get-info", request, response)
	fmt.Println(response)

	// url := fmt.Sprintf("http://127.0.0.1:%s/api/v1/anilist/get-info", config.Port)
	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBodyJSON))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
	// 	return
	// }
	// defer resp.Body.Close()

	// var response schemas.AddAnimeToListResponse
	// if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
	// 	return
	// }

	genres, err := json.Marshal(response.Media[0].Genres)
	if err != nil {
		return
	}

	title := schemas.Title{
		English: response.Media[0].Title.English,
		Romaji:  response.Media[0].Title.Romaji,
	}

	anime := schemas.Anime{
		Aid:          request.Id,
		EpisodeCount: response.Media[0].Episodes,
		Description:  response.Media[0].Description,
		AiringTime:   response.Media[0].NextAiringEpisode.AiringAt,

		InstalledEps: nil,
		Titles:       title,
		Genres:       genres,
	}

	baseService := services.NewBaseService()
	if err := baseService.Add(&anime); err != nil {
		return
	}

	// if err := db.Create(&anime).Error; err != nil {
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": response.Media})
}
