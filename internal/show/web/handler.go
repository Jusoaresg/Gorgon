package web

import (
	"github.com/jmoiron/sqlx"
	tvMazeService "github.com/jusoaresg/gorgon/external/tvmaze/service"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	"github.com/jusoaresg/gorgon/internal/show"
	"github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/internal/show/service"
)

type Handler struct {
	AggregatorService service.ShowAggregatorService
	TvMazeService     tvMazeService.TvMazeSearchService
	ShowManager       tvMazeService.ShowManager
	ShowRepo          repository.ShowRepositoryInterface
	EpisodeRepo       episodeRepository.EpisodeRepositoryInterface
	SeasonRepo        seasonRepository.SeasonRepositoryInterface
	DB                *sqlx.DB
}

func NewHandler(deps *show.Dependencies) *Handler {
	return &Handler{
		AggregatorService: deps.AggregatorService,
		TvMazeService:     deps.TvMazeService,
		ShowManager:       deps.ShowManager,
		ShowRepo:          deps.ShowRepo,
		EpisodeRepo:       deps.EpisodeRepo,
		SeasonRepo:        deps.SeasonRepo,
		DB:                deps.DB,
	}
}
