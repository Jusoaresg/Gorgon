package service

import (
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/pkg/services"
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

	prowlarrHost := configFile.ProwlarrHost
	prowlarrPort := configFile.ProwlarrPort
	return &ProwlarrIndexerService{
		ApiKey:     configFile.ProwlarrApiKey,
		Logger:     logger,
		APIService: *services.NewAPIService(fmt.Sprintf("%s:%s", prowlarrHost, prowlarrPort), logger),
	}
}

func (p *ProwlarrIndexerService) GetIndexers(model *[]schema.IndexerResponse) error {
	return p.APIService.Get(fmt.Sprintf("/api/v1/indexer?apikey=%s", p.ApiKey), &model)
}

func (p *ProwlarrIndexerService) GetIndexer(id int, model *schema.IndexerResponse) error {
	return p.APIService.Get(fmt.Sprintf("/api/v1/indexer/%d?apikey=%s", id, p.ApiKey), &model)
}
