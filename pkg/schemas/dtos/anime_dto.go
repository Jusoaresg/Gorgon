package dtos

type AnilistResponseDto struct {
	Data struct {
		Page struct {
			Media []AnimeDto
		} `json:"Page"`
	} `json:"data"`
}

type AnimeDto struct {
	Id    int
	Title struct {
		English string `json:"english"`
		Romaji  string `json:"romaji"`
	} `json:"title"`
	Description string `json:"description"`

	BannerImage string `json:"bannerImage"`
	CoverImage  struct {
		ExtraLarge string `json:"extraLarge"`
	} `json:"coverImage"`
	Episodes          int      `json:"episodes"`
	Genres            []string `json:"genres"`
	NextAiringEpisode struct {
		AiringAt int `json:"airingAt,omitempty"`
		Episode  int `json:"episode,omitempty"`
	} `json:"nextAiringEpisode"`
	Status string `json:"status"`

	Relations struct {
		Edges []struct {
			Id           int    `json:"id"`
			RelationType string `json:"relationType"`
			Node         struct {
				Id    int `json:"id"`
				Title struct {
					English string `json:"english"`
					Romaji  string `json:"romaji"`
				} `json:"title"`
				Episodes int    `json:"episodes"`
				Format   string `json:"format"`
				Status   string `json:"status"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"relations"`
}
