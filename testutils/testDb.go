package testutils

import (
	"gorgon/migrations"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func GetTestDB() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", ":memory:"+"?_foreign_keys=on")
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	goose.SetBaseFS(migrations.MigrationFS)
	if err := goose.Up(db.DB, "."); err != nil {
		panic(err)
	}

	return db
}
