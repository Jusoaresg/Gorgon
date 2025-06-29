package service

import (
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/pkg/services"
	"log/slog"
	"net/url"
)

type ProwlarrSearchService struct {
	ApiKey     string
	APIService services.APIService
	Logger     *slog.Logger
}

func NewProwlarrSearchService(logger *slog.Logger) (*ProwlarrSearchService, error) {
	configFile, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return &ProwlarrSearchService{
		ApiKey:     configFile.ProwlarrApiKey,
		Logger:     logger,
		APIService: *services.NewAPIService("http://127.0.0.1:9696", logger),
	}, nil
}

func (p *ProwlarrSearchService) Name() string {
	return "ProwlarrSearch"
}

func (p *ProwlarrSearchService) CheckConnection() error {
	var resp struct {
		Status string `json:"status"`
	}
	err := p.APIService.Get("/ping", &resp)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProwlarrSearchService) Search(request *schema.SearchRequest, model *[]schema.SearchResponse) error {
	queryEscaped := url.QueryEscape(request.Query)
	return p.APIService.Get(fmt.Sprintf("/api/v1/search?query=%s&apikey=%s", queryEscaped, p.ApiKey), &model)
}

func (p *ProwlarrSearchService) SearchByType(request *schema.SearchByTypeRequest, model *[]schema.SearchResponse) error {
	queryEscaped := url.QueryEscape(request.Query)
	typeEscaped := url.QueryEscape(request.Type)
	return p.APIService.Get(fmt.Sprintf("/api/v1/search?query=%s&type=%s&apikey=%s", queryEscaped, typeEscaped, p.ApiKey), &model)
}

// func (p *ProwlarrSearchService) GetIndexer(id string, model *schema.IndexerResponse) error {
// 	return p.APIService.Get(fmt.Sprintf("/api/v1/indexer/%s?apikey=%s", id, p.ApiKey), &model)
// }
