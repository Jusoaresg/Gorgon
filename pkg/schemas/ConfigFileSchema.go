package schemas

type ConfigFile struct {
	ProwlarrApiKey string `json:"prowlarrApiKey"`
	ProwlarrHost   string `json:"prowlarrHost"`
	ProwlarrPort   string `json:"prowlarrPort"`

	QBittorrentHost     string `json:"qBittorrentHost"`
	QBittorrentPort     string `json:"qBittorrentPort"`
	QBittorrentUsername string `json:"qBittorrentUsername"`
	QBittorrentPassword string `json:"qBittorrentPassword"`
}
