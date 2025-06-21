package utils

import "os"

func CheckCreateFolder(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(path, 0700); err != nil {
			return err
		}
	}
	return nil
}

func CheckCreateAllFolders(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}
