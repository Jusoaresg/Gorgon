package show

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	tvMazeService "github.com/jusoaresg/gorgon/external/tvmaze/service"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	"github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/internal/show/service"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
)

type Dependencies struct {
	ShowRepo           repository.ShowRepositoryInterface
	ShowAliasesRepo    showAliasRepository.ShowAliasesRepositoryInterface
	AggregatorService  service.ShowAggregatorService
	TvMazeService      tvMazeService.TvMazeSearchService
	ShowManager        tvMazeService.ShowManager
	ShowManagerService *service.ShowManagerService
	EpisodeRepo        episodeRepository.EpisodeRepositoryInterface
	SeasonRepo         seasonRepository.SeasonRepositoryInterface
	DB                 *sqlx.DB
	Logger             *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	tvMazeSvc := tvMazeService.NewTvMazeSearchService(logger)
	showRepo := repository.NewShowRepository(DB)
	showAliasesRepo := showAliasRepository.NewShowAliasesRepository(DB)
	showManagerSvc := service.NewShowManagerService(logger, DB)

	return &Dependencies{
		ShowRepo:           showRepo,
		ShowAliasesRepo:    &showAliasesRepo,
		AggregatorService:  *service.NewShowAggregatorServiceWithDb(DB),
		TvMazeService:      *tvMazeSvc,
		ShowManager:        *tvMazeService.NewShowManager(*tvMazeSvc, *showRepo, logger),
		ShowManagerService: showManagerSvc,
		EpisodeRepo:        episodeRepository.NewEpisodeRepository(DB),
		SeasonRepo:         seasonRepository.NewSeasonRepository(DB),
		DB:                 DB,
		Logger:             logger,
	}
}
