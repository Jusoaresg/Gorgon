package config

import (
	"github.com/jusoaresg/gorgon/utils"
)

func InitializeDownloadFolders() error {

	downloadFolder := "downloads"

	if err := utils.CheckCreateFolder(downloadFolder); err != nil {
		return err
	}

	return nil
}
