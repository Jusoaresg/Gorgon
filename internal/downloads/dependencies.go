package downloads

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/internal/downloads/service"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	episodeTorrentRepository "github.com/jusoaresg/gorgon/internal/episode_torrent/repository"
)

type Dependencies struct {
	Service            *service.DownloadsService
	EpisodeRepo        episodeRepository.EpisodeRepositoryInterface
	EpisodeTorrentRepo episodeTorrentRepository.EpisodeTorrentRepositoryInterface
	DB                 *sqlx.DB
	Logger             *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	return &Dependencies{
		Service:            service.NewDownloadsService(DB, logger),
		EpisodeRepo:        episodeRepository.NewEpisodeRepository(DB),
		EpisodeTorrentRepo: episodeTorrentRepository.NewEpisodeTorrentRepository(DB),
		DB:                 DB,
		Logger:             logger,
	}
}
