package web

import (
	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/internal/downloads"
	"github.com/jusoaresg/gorgon/internal/downloads/service"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	episodeTorrentRepository "github.com/jusoaresg/gorgon/internal/episode_torrent/repository"
)

type Handler struct {
	Service            *service.DownloadsService
	EpisodeRepo        episodeRepository.EpisodeRepositoryInterface
	EpisodeTorrentRepo episodeTorrentRepository.EpisodeTorrentRepositoryInterface
	DB                 *sqlx.DB
}

func NewHandler(deps *downloads.Dependencies) *Handler {
	return &Handler{
		Service:            deps.Service,
		EpisodeRepo:        deps.EpisodeRepo,
		EpisodeTorrentRepo: deps.EpisodeTorrentRepo,
		DB:                 deps.DB,
	}
}
