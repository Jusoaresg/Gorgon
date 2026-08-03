package schema

type AddEpisodeTorrentRequest struct {
	EpisodeID   int64  `json:"episodeID"`
	InfoHash    string `json:"infoHash"`
	MagneticUrl string `json:"magneticUrl"`

	Title       string `json:"title,omitempty"`
	Indexer     string `json:"indexer,omitempty"`
	InfoUrl     string `json:"infoUrl,omitempty"`
	PublishDate string `json:"publishDate,omitempty"`
}
