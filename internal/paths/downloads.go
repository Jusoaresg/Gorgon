package paths

import (
	"github.com/jusoaresg/gorgon/config"
	"path/filepath"
)

func GetTorrentDownloadFolder() (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	return filepath.Join(cfg.QBittorrentDownloadFolder), nil
}

func GetEpisodeDownloadFile(episodeFileName string) (string, error) {
	downloadFolder, err := GetTorrentDownloadFolder()
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filepath.Join(downloadFolder, episodeFileName))
	if err != nil {
		return "", err
	}
	return absPath, nil
}
