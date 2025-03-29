package schemas

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

type IndexerResponse struct {
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	Enabled        bool     `json:"enable"`
	DefinitionName string   `json:"definitionName"`
	IndexerUrls    []string `json:"indexerUrls"`
	Language       string   `json:"language"`
}

type RemoveIndexerRequest struct {
	Id string `json:"id"`
}
