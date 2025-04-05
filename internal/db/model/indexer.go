package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Indexer struct {
	gorm.Model
	IndexerId      int            `json:"id"`
	Name           string         `json:"name"`
	Enabled        bool           `json:"enable"`
	DefinitionName string         `json:"definitionName"`
	IndexerUrls    datatypes.JSON `json:"indexerUrls"`
	Language       string         `json:"language"`
}
