package model

type Indexer struct {
	ID             int64  `db:"id" json:"internalId"`
	IndexerID      int    `db:"indexer_id" json:"id"`
	Name           string `db:"name" json:"name"`
	Enabled        bool   `db:"enabled" json:"enabled"`
	DefinitionName string `db:"definition_name" json:"definitionName"`
	IndexerUrls    string `db:"indexer_urls" json:"indexerUrls"`
	Language       string `db:"language" json:"language"`
}
