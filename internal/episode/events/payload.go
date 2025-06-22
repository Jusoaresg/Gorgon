package episode

type EpisodeTrackingUpdatedPayload struct {
	Type      string `json:"type"`
	EpisodeID int64  `json:"episodeID"`
	Tracking  string `json:"tracking"`
}
