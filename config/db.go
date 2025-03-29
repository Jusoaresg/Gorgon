package config

import (
	"gorgon/pkg/schemas"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDb() (*gorm.DB, error) {

	dbPath := "assets/anime.db"

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
		&schemas.Anime{},
		&schemas.Title{},
		&schemas.Indexer{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
