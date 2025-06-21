package config

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

var (
	ConfigPath  string = "assets/config.json"
	Port        string = "8080"
	baseApiPath string
	db          *sqlx.DB
	logger      *slog.Logger
)

func Init() error {
	var err error

	dbInstance, err := InitializeDb()
	if err != nil {
		return err
	}
	db = dbInstance

	InitializeOrUpdateConfigFile()

	InitializeDownloadFolders()

	return nil
}

func GetLogger() *slog.Logger {
	return NewLogger()
}

func GetSQLite() *sqlx.DB {
	if db == nil {
		panic("database is not initialized")
	}
	return db
}
