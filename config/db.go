package config

import (
	"fmt"
	"os"

	"gorgon/internal/db/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDb() (*gorm.DB, error) {

	dbFolder := "assets"
	dbPath := fmt.Sprintf("%s/%s", dbFolder, "gorgon.db")

	_, err := os.Stat(dbFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(dbFolder, 0700)
		if err != nil {
			return nil, err
		}
	}

	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.Show{},
		&model.Episode{},
		&model.EpisodeContent{},
		&model.Schedule{},
		&model.Externals{},
		&model.Image{},
		&model.Season{},

		&model.Indexer{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
