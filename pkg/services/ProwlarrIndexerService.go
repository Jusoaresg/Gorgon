package services

import (
	"fmt"
	"gorgon/config"
	"gorgon/pkg/schemas"
)

type ProwlarrIndexerService struct {
	ApiKey     string
	APIService APIService
}

func NewProwlarrIndexerService() *ProwlarrIndexerService {
	configFile, err := config.LoadConfig()
	if err != nil {
		panic("Error while creting Prowlarr Indexer Service")
	}

	return &ProwlarrIndexerService{
		ApiKey:     configFile.ProwlarrApiKey,
		APIService: *NewAPIService("http://127.0.0.1:9696"),
	}
}

func (p *ProwlarrIndexerService) GetIndexers(model *[]schemas.IndexerResponse) error {
	return p.APIService.Get(fmt.Sprintf("/api/v1/indexer?apikey=%s", p.ApiKey), &model)
}

func (p *ProwlarrIndexerService) GetIndexer(id string, model *schemas.IndexerResponse) error {
	return p.APIService.Get(fmt.Sprintf("/api/v1/indexer/%s?apikey=%s", id, p.ApiKey), &model)
}
