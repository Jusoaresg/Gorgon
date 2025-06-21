package utils

import (
	"os"
	"path/filepath"
)

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func DeleteFile(basePath, fileName string) (string, error) {
	filePath := filepath.Join(basePath, fileName)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	if !FileExists(absPath) {
		return "", os.ErrNotExist
	}

	err = os.Remove(absPath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
