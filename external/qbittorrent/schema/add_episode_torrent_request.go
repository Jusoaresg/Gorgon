package schema

type AddEpisodeTorrentRequest struct {
	EpisodeID   int64  `json:"episodeID"`
	InfoHash    string `json:"infoHash"`
	MagneticUrl string `json:"magneticUrl"`
}
