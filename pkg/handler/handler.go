package handler

import (
	"github.com/jusoaresg/gorgon/config"

	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func InitHandler() {
	db = config.GetSQLite()
}
