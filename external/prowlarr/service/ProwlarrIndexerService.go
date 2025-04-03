package service

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/prowlarr/schema"
	"gorgon/pkg/services"
	"log/slog"
)

type ProwlarrIndexerService struct {
	ApiKey     string
	APIService services.APIService
	Logger     *slog.Logger
}

func NewProwlarrIndexerService(logger *slog.Logger) *ProwlarrIndexerService {
	configFile, err := config.LoadConfig()
	if err != nil {
		panic("Error while creting Prowlarr Indexer Service")
	}

	return &ProwlarrIndexerService{
		ApiKey:     configFile.ProwlarrApiKey,
		Logger:     logger,
		APIService: *services.NewAPIService("http://127.0.0.1:9696", logger),
	}
}

func (p *ProwlarrIndexerService) GetIndexers(model *[]schema.IndexerResponse) error {
	return p.APIService.Get(fmt.Sprintf("/api/v1/indexer?apikey=%s", p.ApiKey), &model)
}

func (p *ProwlarrIndexerService) GetIndexer(id int, model *schema.IndexerResponse) error {
	return p.APIService.Get(fmt.Sprintf("/api/v1/indexer/%d?apikey=%s", id, p.ApiKey), &model)
}
