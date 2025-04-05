package config

import (
	"os"

	"gorgon/internal/db/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDb() (*gorm.DB, error) {

	dbPath := "assets/gorgon.db"

	_, err := os.Stat(dbPath)
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

	err = db.AutoMigrate(
		&model.Show{},
		&model.Episode{},
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
