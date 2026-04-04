package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/pkg/services"
	"golang.org/x/time/rate"
)

type ErrProwlarrHostPortNotSet interface {
	Error() string
	IsProwlarrConfigWarn() bool
}
type ProwlarrHostPortNotSet struct{}

func (m *ProwlarrHostPortNotSet) Error() string {
	return "Prowlar Host or Port not set"
}

func (m *ProwlarrHostPortNotSet) IsProwlarrConfigWarn() bool {
	return true
}

type ProwlarrSearchService struct {
	ApiKey     string
	APIService services.APIService
	Logger     *slog.Logger
	limiter    *rate.Limiter
}

func NewProwlarrSearchService(logger *slog.Logger) (*ProwlarrSearchService, error) {
	configFile, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	prowlarrHost := configFile.ProwlarrHost
	prowlarrPort := configFile.ProwlarrPort
	return &ProwlarrSearchService{
		ApiKey:     configFile.ProwlarrApiKey,
		Logger:     logger,
		APIService: *services.NewAPIService(fmt.Sprintf("%s:%s", prowlarrHost, prowlarrPort), logger),
		limiter:    rate.NewLimiter(rate.Every(2*time.Second), 1),
	}, nil
}

func (p *ProwlarrSearchService) Name() string {
	return "ProwlarrSearch"
}

func (p *ProwlarrSearchService) waitTicker(ctx context.Context) error {
	if err := p.limiter.Wait(ctx); err != nil {
		return err
	}
	return nil
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
	if err := p.waitTicker(context.Background()); err != nil {
		return err
	}
	queryEscaped := url.QueryEscape(request.Query)
	return p.APIService.Get(fmt.Sprintf("/api/v1/search?query=%s&apikey=%s", queryEscaped, p.ApiKey), &model)
}

func (p *ProwlarrSearchService) SearchByType(request *schema.SearchByTypeRequest, model *[]schema.SearchResponse) error {
	if err := p.waitTicker(context.Background()); err != nil {
		return err
	}
	queryEscaped := url.QueryEscape(request.Query)
	typeEscaped := url.QueryEscape(request.Type)
	return p.APIService.Get(fmt.Sprintf("/api/v1/search?query=%s&type=%s&apikey=%s", queryEscaped, typeEscaped, p.ApiKey), &model)
}
