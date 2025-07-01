package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

var (
	InDocker            = false
	ConfigFolder string = "./configs"

	LogsPath    string = filepath.Join(ConfigFolder, "logs")
	Port        string = "8080"
	baseApiPath string
	db          *sqlx.DB
	logger      *slog.Logger
)

func Init() error {
	var err error
	InDocker = os.Getenv("IN_DOCKER") == "true"

	if InDocker {
		ConfigFolder = "/configs"
	}

	LogsPath = filepath.Join(ConfigFolder, "logs")
	logger = NewLogger()

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
	if logger == nil {
		logger = NewLogger()
	}
	return logger
}

func GetSQLite() *sqlx.DB {
	if db == nil {
		panic("database is not initialized")
	}
	return db
}
