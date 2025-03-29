package schemas

type AnimeGetInfoByIdRequest struct {
	Id string
}

type AnimeGetInfoByIdResponse struct {
	Data struct {
		Page struct {
			Media []struct {
				Title struct {
					English string `json:"english"`
					Romaji  string `json:"romaji"`
				}

				NextAiringEpisode struct {
					Episode  interface{}
					AiringAt interface{}
				}

				Episodes    int
				Description string
				Genres      []string
			}
		}
	}
}
