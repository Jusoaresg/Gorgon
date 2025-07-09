package config

import (
	"os"
	"path/filepath"

	"github.com/jusoaresg/gorgon/migrations"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func InitializeDb() (*sqlx.DB, error) {

	dbPath := filepath.Join(ConfigFolder, "gorgon.db")

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(ConfigFolder, 0700)
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

	db, err := sqlx.Open("sqlite", "file:"+dbPath+"?_pragma=foreign_keys=ON&_pragma=journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(migrations.MigrationFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	if err := goose.Up(db.DB, "."); err != nil {
		panic(err)
	}

	return db, nil
}
