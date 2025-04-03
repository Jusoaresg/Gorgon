package service

import (
	"gorgon/external/anilist/schema"
	"gorgon/pkg/schemas/dtos"
	"gorgon/pkg/services"
	"log/slog"
)

type AnilistService struct {
	APIService *services.APIService
	Logger     *slog.Logger
}

func NewAnimeListService(logger *slog.Logger) *AnilistService {
	return &AnilistService{
		APIService: services.NewAPIService("https://graphql.anilist.co", logger),
		Logger:     logger,
	}
}

// func (a *AnilistService) GetAnimeInfoById(id int) (*schema.AnimeGetInfoByIdResponse, error) {
func (a *AnilistService) GetAnimeInfoById(id int) (*dtos.AnimeDto, error) {
	query := `
		query($mediaId: Int)  {
		  Page {
		    media(type: ANIME, id: $mediaId) {
		      id
		      title {
			romaji
			english
		      }
		      nextAiringEpisode {
						episode
						airingAt
					}
		      description
		      episodes
		      genres
		      bannerImage
		      coverImage {
			extraLarge
		      }
		      relations {
			edges {
				id
				relationType
				node {
					id
					title {
						english
						romaji
					}
					status
					format
					episodes
				}
			}
		      }
		      status
		    }
		    }
		  }
	`

	var variables = map[string]interface{}{
		"mediaId": id,
	}

	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	var response dtos.AnilistResponseDto
	if _, err := a.APIService.Post("", &reqBody, &response); err != nil {
		a.Logger.Info("Failed to fetch anime info by ID from Anilist", slog.String("error", err.Error()), slog.Int("Id", id))
		return nil, err
	}

	return &response.Data.Page.Media[0], nil
}

func (a *AnilistService) FindAnimeIdAnilist(request *schema.AnimeTitleIdByNameRequest) (*schema.AnimeTitleIdByNameResponse, error) {
	query := ` 
		query ($search: String!) {
			Page {
				media(search: $search, type: ANIME) {
					id
					title {
						romaji
						english
					}
				}
			}
		}
	`

	var variables = map[string]interface{}{
		"search": request.Name,
	}

	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	var response schema.AnimeTitleIdByNameResponse
	if _, err := a.APIService.Post("", &reqBody, &response); err != nil {
		a.Logger.Error("Failed to search for anime by name from anilist", slog.String("error", err.Error()))
		return nil, err
	}

	return &response, nil
}
