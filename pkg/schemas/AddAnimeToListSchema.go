package schemas

type AddAnimeToListRequest struct {
	Id string `json:"id"`
}

type AddAnimeToListResponse struct {
	Media []struct {
		Title struct {
			English string `json:"english"`
			Romaji  string `json:"romaji"`
		}

		NextAiringEpisode struct {
			Episode  string
			AiringAt string
		}

		Episodes    int
		Description string
		Genres      []string
	} `json:"media"`
}
