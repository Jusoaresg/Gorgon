package episode

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	"github.com/jusoaresg/gorgon/internal/episode/service"
	epContentRepository "github.com/jusoaresg/gorgon/internal/episode_content/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
)

type Dependencies struct {
	EpisodeRepo        episodeRepository.EpisodeRepositoryInterface
	EpisodeContentRepo epContentRepository.EpisodeContentRepository
	ShowRepo           showRepository.ShowRepositoryInterface
	EpisodeSearchSvc   *service.EpisodeSearchService
	DB                 *sqlx.DB
	Logger             *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	episodeRepo := episodeRepository.NewEpisodeRepository(DB)
	showRepo := showRepository.NewShowRepository(DB)

	return &Dependencies{
		EpisodeRepo:        episodeRepo,
		EpisodeContentRepo: epContentRepository.NewEpisodeContentRepository(DB),
		ShowRepo:           showRepo,
		EpisodeSearchSvc: service.NewEpisodeSearchService(
			DB,
			logger,
			&service.EpisodeSearcher{},
			&service.EpisodeDownloader{},
			episodeRepo,
			showRepo,
		),
		DB:     DB,
		Logger: logger,
	}
}
