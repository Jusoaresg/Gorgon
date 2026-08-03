package episode

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	"github.com/jusoaresg/gorgon/internal/episode/service"
	epContentRepository "github.com/jusoaresg/gorgon/internal/episode_content/repository"
	episodeTorrentRepository "github.com/jusoaresg/gorgon/internal/episode_torrent/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showAliasesRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
)

type Dependencies struct {
	EpisodeRepo        episodeRepository.EpisodeRepositoryInterface
	EpisodeContentRepo epContentRepository.EpisodeContentRepository
	EpisodeTorrentRepo episodeTorrentRepository.EpisodeTorrentRepositoryInterface
	ShowRepo           showRepository.ShowRepositoryInterface
	ShowAliasesRepo    showAliasesRepository.ShowAliasesRepositoryInterface
	EpisodeSearchSvc   *service.EpisodeSearchService
	DB                 *sqlx.DB
	Logger             *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	episodeRepo := episodeRepository.NewEpisodeRepository(DB)
	showRepo := showRepository.NewShowRepository(DB)
	showAliasesRepo := showAliasesRepository.NewShowAliasesRepository(DB)

	return &Dependencies{
		EpisodeRepo:        episodeRepo,
		EpisodeContentRepo: epContentRepository.NewEpisodeContentRepository(DB),
		EpisodeTorrentRepo: episodeTorrentRepository.NewEpisodeTorrentRepository(DB),
		ShowRepo:           showRepo,
		ShowAliasesRepo:    &showAliasesRepo,
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
