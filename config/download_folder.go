package config

import (
	"github.com/jusoaresg/gorgon/utils"
	"path/filepath"
)

func InitializeDownloadFolders() error {

	assetsPath := "assets/"

	downloadFolder := filepath.Join(assetsPath, "downloads")

	if err := utils.CheckCreateFolder(downloadFolder); err != nil {
		return err
	}

	return nil
}
