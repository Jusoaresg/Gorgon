package utils

import (
	"fmt"
	"os"
	"path/filepath"

	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	epContentModel "github.com/jusoaresg/gorgon/internal/episode_content/model"
)

func SymlinkPathForEpisode(showsFolder, showName string, episode episodeModel.Episode, episodeContent epContentModel.EpisodeContent) (string, error) {
	seasonFolder := fmt.Sprintf("Season %d", episode.Season)
	destFolder := filepath.Join(showsFolder, showName, seasonFolder)
	if err := CheckCreateAllFolders(destFolder); err != nil {
		return "", err
	}

	fileExtension := filepath.Ext(episodeContent.Name)
	episodeNewName := fmt.Sprintf("%s - S%dE%d - %s%s", showName, episode.Season, episode.Number, episode.Name, fileExtension)
	destPath, err := filepath.Abs(filepath.Join(destFolder, episodeNewName))
	if err != nil {
		return "", err
	}
	return destPath, nil
}

func DeleteSymlink(showsFolder, showName string, episode episodeModel.Episode, episodeContent epContentModel.EpisodeContent) error {
	linkPath, err := SymlinkPathForEpisode(showsFolder, showName, episode, episodeContent)
	if err != nil {
		return err
	}

	file, err := os.Lstat(linkPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if file.Mode()&os.ModeSymlink != 0 {
		return os.Remove(linkPath)
	}

	return nil
}

func CreateSymlink(downloadPath, destPath string) error {
	if _, err := os.Lstat(destPath); err == nil {
		if err := os.Remove(destPath); err != nil {
			return err
		}
	}

	relativePath, err := filepath.Rel(filepath.Dir(destPath), downloadPath)
	if err != nil {
		return err
	}

	return os.Symlink(relativePath, destPath)
}

func IsSymlinkBroken(path string) (bool, error) {
	info, err := os.Lstat(
		path,
	)
	if err != nil {
		return false, err
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return false, nil
	}
	target, err := os.Readlink(path)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(target)
	return os.IsNotExist(err), nil
}
