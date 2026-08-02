package episode

type EpisodeTrackingUpdatedPayload struct {
	Type      string `json:"type"`
	EpisodeID int64  `json:"episodeID"`
	Tracking  string `json:"tracking"`
	InfoUrl   string `json:"infoUrl"`
}

type EpisodeSearchFinishedPayload struct {
	Type      string `json:"type"`
	EpisodeID int64  `json:"episodeID"`
	Season    int    `json:"season"`
	Number    int    `json:"number"`
	Name      string `json:"name"`
	ShowName  string `json:"showName"`
	Result    string `json:"result"`
	Message   string `json:"message"`
}
