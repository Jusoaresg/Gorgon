package schema

type ChangeEpisodeTrackingRequest struct {
	EpisodeIds []int  `json:"episode_ids"`
	Tracking   string `json:"tracking"`
}
