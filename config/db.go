package config

import (
	"fmt"
	"github.com/jusoaresg/gorgon/migrations"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func InitializeDb() (*sqlx.DB, error) {

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

	db, err := sqlx.Open("sqlite3", dbPath+"?_foreign_keys=on")
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
