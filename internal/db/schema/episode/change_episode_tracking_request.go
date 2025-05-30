package episode

type ChangeEpisodeTrackingRequest struct {
	EpisodeId int    `json:"episode_id"`
	Tracking  string `json:"tracking"`
}
