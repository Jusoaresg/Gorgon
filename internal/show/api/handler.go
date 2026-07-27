package api

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	tvMazeService "github.com/jusoaresg/gorgon/external/tvmaze/service"
	"github.com/jusoaresg/gorgon/internal/show"
	"github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/internal/show/service"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
)

type Handler struct {
	ShowRepo           repository.ShowRepositoryInterface
	ShowAliasRepo      showAliasRepository.ShowAliasesRepositoryInterface
	AggregatorService  service.ShowAggregatorService
	TvMazeService      tvMazeService.TvMazeSearchService
	ShowManager        tvMazeService.ShowManager
	ShowManagerService *service.ShowManagerService
	DB                 *sqlx.DB
	Logger             *slog.Logger
}

func NewHandler(deps *show.Dependencies) *Handler {
	return &Handler{
		ShowRepo:           deps.ShowRepo,
		ShowAliasRepo:      deps.ShowAliasRepo,
		AggregatorService:  deps.AggregatorService,
		TvMazeService:      deps.TvMazeService,
		ShowManager:        deps.ShowManager,
		ShowManagerService: deps.ShowManagerService,
		DB:                 deps.DB,
		Logger:             deps.Logger,
	}
}
