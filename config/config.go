package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/jmoiron/sqlx"
)

type SafeDB struct {
	Db    *sqlx.DB
	Write *sync.Mutex
}

var (
	InDocker            = false
	BaseDir             = "."
	ConfigFolder string = "./configs"

	LogsPath string = filepath.Join(ConfigFolder, "logs")
	Port     string = "8080"
	safeDB   SafeDB
	logger   *slog.Logger
)

func Init() error {
	InDocker = os.Getenv("IN_DOCKER") == "true"

	if baseDirEnv := os.Getenv("GORGON_BASE_DIR"); baseDirEnv != "" {
		BaseDir = baseDirEnv
	}

	if InDocker {
		ConfigFolder = "/configs"
		LogsPath = filepath.Join(ConfigFolder, "logs")
		BaseDir = "/"
	}

	if err := ReloadFolders(); err != nil {
		return err
	}

	logger = NewLogger()

	dbInstance, err := InitializeDb()
	if err != nil {
		return err
	}
	safeDB = SafeDB{
		Db:    dbInstance,
		Write: &sync.Mutex{},
	}

	InitializeOrUpdateConfigFile()

	return nil
}

func ReloadFolders() error {
	ConfigFolder = filepath.Join(BaseDir, "configs")
	LogsPath = filepath.Join(ConfigFolder, "logs")

	dirs := []string{
		ConfigFolder,
		LogsPath,
		filepath.Join(BaseDir, "downloads"),
	}

	for _, dir := range dirs {
		//TODO: Logging creating necessary directories
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func GetLogger() *slog.Logger {
	if logger == nil {
		logger = NewLogger()
	}
	return logger
}

func GetSQLite() *sqlx.DB {
	if safeDB.Db == nil {
		panic("database is not initialized")
	}
	return safeDB.Db
}

func GetSafeDB() *SafeDB {
	if safeDB.Db == nil {
		panic("database is not initialized")
	}
	return &safeDB
}
