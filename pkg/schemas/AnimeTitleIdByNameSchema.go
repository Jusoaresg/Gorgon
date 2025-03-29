package schemas

type AnimeTitleIdByNameRequest struct {
	Name string `json:"name"`
}

type AnimeTitleIdByNameResponse struct {
	Data struct {
		Page struct {
			Media []struct {
				ID    int `json:"id"`
				Title struct {
					Romaji  string `json:"romaji"`
					English string `json:"english"`
				} `json:"title"`
			} `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}
