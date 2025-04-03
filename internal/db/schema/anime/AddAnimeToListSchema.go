package anime

type AddAnimeToListResponse struct {
	Data []struct {
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
		BannerImage string
		CoverImage  struct {
			ExtraLarge string
		}
	} `json:"data"`
}
