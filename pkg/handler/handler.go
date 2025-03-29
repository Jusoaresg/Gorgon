package handler

import (
	"gorgon/config"

	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func InitHandler() {
	db = config.GetSQLite()
}
