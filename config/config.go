package config

import (
	"gorm.io/gorm"
)

var (
	ConfigPath      string = "assets/config.json"
	Port            string = "8080"
	baseApiPath     string
	animeIdDataPath string = "assets/anime-titles.dat.gz"
	db              *gorm.DB
)

func Init() error {
	var err error

	db, err = InitializeDb()
	if err != nil {
		return err
	}

	InitializeConfigFile()

	return nil
}

func GetSQLite() *gorm.DB {
	return db
}

func GetAnimeIdDataPath() string {
	return animeIdDataPath
}
