package schemas

type ConfigFile struct {
	ProwlarrApiKey string `json:"prowlarrApiKey"`
	ProwlarrHost   string `json:"prowlarrHost"`
	ProwlarrPort   string `json:"prowlarrPort"`

	QBittorrentHost           string `json:"qBittorrentHost"`
	QBittorrentPort           string `json:"qBittorrentPort"`
	QBittorrentUsername       string `json:"qBittorrentUsername"`
	QBittorrentPassword       string `json:"qBittorrentPassword"`
	QBittorrentDownloadFolder string `json:"qBittorrentDownloadFolder"`

	DefaultShowInfoFolder string `json:"defaultShowInfoFolder"`
	ShowsFolder           string `json:"showsFolder"`
}

// NOTE: For patch route
type UpdateConfigInput struct {
	ProwlarrApiKey *string `json:"prowlarrApiKey"`
	ProwlarrHost   *string `json:"prowlarrHost"`
	ProwlarrPort   *string `json:"prowlarrPort"`

	QBittorrentHost           *string `json:"qBittorrentHost"`
	QBittorrentPort           *string `json:"qBittorrentPort"`
	QBittorrentUsername       *string `json:"qBittorrentUsername"`
	QBittorrentPassword       *string `json:"qBittorrentPassword"`
	QBittorrentDownloadFolder *string `json:"qBittorrentDownloadFolder"`

	DefaultShowInfoFolder *string `json:"defaultShowInfoFolder"`
	ShowsFolder           *string `json:"showsFolder"`
}

func setString(dst *string, src *string) {
	if src != nil {
		*dst = *src
	}
}

func (c *ConfigFile) Apply(input *UpdateConfigInput) {
	setString(&c.ProwlarrApiKey, input.ProwlarrApiKey)
	setString(&c.ProwlarrHost, input.ProwlarrHost)
	setString(&c.ProwlarrPort, input.ProwlarrPort)

	setString(&c.QBittorrentHost, input.QBittorrentHost)
	setString(&c.QBittorrentPort, input.QBittorrentPort)
	setString(&c.QBittorrentUsername, input.QBittorrentUsername)
	setString(&c.QBittorrentPassword, input.QBittorrentPassword)
	setString(&c.QBittorrentDownloadFolder, input.QBittorrentDownloadFolder)

	setString(&c.DefaultShowInfoFolder, input.DefaultShowInfoFolder)
	setString(&c.ShowsFolder, input.ShowsFolder)
}
