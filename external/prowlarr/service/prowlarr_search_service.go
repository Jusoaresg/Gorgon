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

var (
	prowlarrSemaphore   = make(chan struct{}, 3)
	prowlarrRateLimiter = rate.NewLimiter(rate.Every(2*time.Second), 1)
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

	host := configFile.ProwlarrHost
	port := configFile.ProwlarrPort

	if host == "" || port == "" {
		return nil, &ProwlarrHostPortNotSet{}
	}

	return &ProwlarrSearchService{
		ApiKey:     configFile.ProwlarrApiKey,
		Logger:     logger,
		APIService: *services.NewAPIService(fmt.Sprintf("%s:%s", host, port), logger),
	}, nil
}

func (p *ProwlarrSearchService) Name() string {
	return "ProwlarrSearch"
}

func (p *ProwlarrSearchService) waitAndLock(ctx context.Context) error {
	if err := prowlarrRateLimiter.Wait(ctx); err != nil {
		return err
	}
	prowlarrSemaphore <- struct{}{}
	return nil
}

func (p *ProwlarrSearchService) unlock() {
	<-prowlarrSemaphore
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
	if err := p.waitAndLock(context.Background()); err != nil {
		return err
	}
	defer p.unlock()

	queryEscaped := url.QueryEscape(request.Query)
	return p.APIService.Get(fmt.Sprintf("/api/v1/search?query=%s&apikey=%s", queryEscaped, p.ApiKey), &model)
}

func (p *ProwlarrSearchService) SearchByType(request *schema.SearchByTypeRequest, model *[]schema.SearchResponse) error {
	if err := p.waitAndLock(context.Background()); err != nil {
		return err
	}
	defer p.unlock()

	queryEscaped := url.QueryEscape(request.Query)
	typeEscaped := url.QueryEscape(request.Type)
	return p.APIService.Get(fmt.Sprintf("/api/v1/search?query=%s&type=%s&apikey=%s", queryEscaped, typeEscaped, p.ApiKey), &model)
}
