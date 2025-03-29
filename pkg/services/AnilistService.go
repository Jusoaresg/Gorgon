package services

import (
	"gorgon/pkg/schemas"
)

type AnilistService struct {
	APIService *APIService
}

func NewAnimeListService() *AnilistService {
	return &AnilistService{
		APIService: NewAPIService("https://graphql.anilist.co"),
	}
}

func (a *AnilistService) GetAnimeInfoById(request *schemas.AnimeGetInfoByIdRequest) *schemas.AnimeGetInfoByIdResponse {
	query := `
		query ($mediaId: Int) {
		  Page {
		    media(id: $mediaId, type: ANIME) {
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
		    }
		  }
		}
`

	var variables = map[string]interface{}{
		"mediaId": request.Id,
	}

	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	var response schemas.AnimeGetInfoByIdResponse
	if err := a.APIService.Post("", &reqBody, &response); err != nil {
		return nil
	}

	return &response
}

func (a *AnilistService) FindAnimeIdAnilist(request *schemas.AnimeTitleIdByNameRequest) *schemas.AnimeTitleIdByNameResponse {
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

	var response schemas.AnimeTitleIdByNameResponse
	if err := a.APIService.Post("", &reqBody, &response); err != nil {
		return nil
	}

	return &response
}
