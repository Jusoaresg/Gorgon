package schema

type IndexerResponse struct {
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	Enabled        bool     `json:"enable"`
	DefinitionName string   `json:"definitionName"`
	IndexerUrls    []string `json:"indexerUrls"`
	Language       string   `json:"language"`
}
