package dtos

type ShowDto struct {
	ShowID    int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Language  string `json:"language"`
	Status    string `json:"status"`
	Premiered string `json:"premiered"`
	Ended     string `json:"ended"`
	Rating    struct {
		Average *float64 `json:"average"`
	} `json:"rating,omitempty"`
	Summary string `json:"summary"`
	Updated int    `json:"updated"`

	Schedule struct {
		Days []string `json:"days"`
		Time string   `json:"time"`
	} `json:"schedule"`
	Externals struct {
		Imdb    string `json:"imdb"`
		TheTvdb int    `json:"thetvdb"`
		TvRage  int    `json:"tvrage"`
	} `json:"externals"`
	Image struct {
		Medium   string `json:"medium"`
		Original string `json:"original"`
	} `json:"image"`
	Genres []string `json:"genres"`
}

// Always use as slice
type SeasonDto struct {
	ShowId       int    `json:"id"`
	Url          string `json:"url"`
	Name         string `json:"name,omitempty"`
	Number       int    `json:"number"`
	PremiereDate string `json:"premiereDate"`
	EndDate      string `json:"endDate"`
}

// Always use as slice
type EpisodeDto struct {
	ShowId   int    `json:"id"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	Season   int    `json:"season"`
	Number   int    `json:"number"`
	AirDate  string `json:"airdate"`
	AirStamp string `json:"airStamp"`
	AirTime  string `json:"airTime"`
	Summary  string `json:"summary"`
}
