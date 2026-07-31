package web

import (
	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/internal/episode"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	"github.com/jusoaresg/gorgon/internal/show/repository"
	showAliasesRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
)

type Handler struct {
	ShowRepo        repository.ShowRepositoryInterface
	EpisodeRepo     episodeRepository.EpisodeRepositoryInterface
	ShowAliasesRepo showAliasesRepository.ShowAliasesRepositoryInterface
	DB              *sqlx.DB
}

func NewHandler(deps *episode.Dependencies) *Handler {
	return &Handler{
		ShowRepo:        deps.ShowRepo,
		ShowAliasesRepo: deps.ShowAliasesRepo,
		EpisodeRepo:     deps.EpisodeRepo,
		DB:              deps.DB,
	}
}
