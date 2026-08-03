package web

import (
	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/internal/episode"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	episodeTorrentRepository "github.com/jusoaresg/gorgon/internal/episode_torrent/repository"
	"github.com/jusoaresg/gorgon/internal/show/repository"
	showAliasesRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
)

type Handler struct {
	ShowRepo           repository.ShowRepositoryInterface
	EpisodeRepo        episodeRepository.EpisodeRepositoryInterface
	EpisodeTorrentRepo episodeTorrentRepository.EpisodeTorrentRepositoryInterface
	ShowAliasesRepo    showAliasesRepository.ShowAliasesRepositoryInterface
	DB                 *sqlx.DB
}

func NewHandler(deps *episode.Dependencies) *Handler {
	return &Handler{
		ShowRepo:           deps.ShowRepo,
		ShowAliasesRepo:    deps.ShowAliasesRepo,
		EpisodeRepo:        deps.EpisodeRepo,
		EpisodeTorrentRepo: deps.EpisodeTorrentRepo,
		DB:                 deps.DB,
	}
}
