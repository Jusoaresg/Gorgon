package config

import (
	"log/slog"

	"gorm.io/gorm"
)

var (
	ConfigPath  string = "assets/config.json"
	Port        string = "8080"
	baseApiPath string
	db          *gorm.DB
	logger      *slog.Logger
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

func GetLogger() *slog.Logger {
	return NewLogger()
}

func GetSQLite() *gorm.DB {
	return db
}
