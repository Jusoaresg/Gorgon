package service

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/prowlarr/schema"
	"gorgon/pkg/services"
	"log/slog"
)

type ProwlarrSearchService struct {
	ApiKey     string
	APIService services.APIService
	Logger     *slog.Logger
}

func NewProwlarrSearchService(logger *slog.Logger) *ProwlarrSearchService {
	configFile, err := config.LoadConfig()
	if err != nil {
		panic("Error while creting Prowlarr Query Service")
	}

	return &ProwlarrSearchService{
		ApiKey:     configFile.ProwlarrApiKey,
		Logger:     logger,
		APIService: *services.NewAPIService("http://127.0.0.1:9696", logger),
	}
}

func (p *ProwlarrSearchService) Search(request *schema.SearchRequest, model *[]schema.SearchResponse) error {
	return p.APIService.Get(fmt.Sprintf("/api/v1/search?query='%s'&apikey=%s", request.Query, p.ApiKey), &model)
}

// func (p *ProwlarrSearchService) GetIndexer(id string, model *schema.IndexerResponse) error {
// 	return p.APIService.Get(fmt.Sprintf("/api/v1/indexer/%s?apikey=%s", id, p.ApiKey), &model)
// }
