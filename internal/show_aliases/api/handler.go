package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/internal/show_aliases"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
)

type Handler struct {
	ShowAliasRepo showAliasRepository.ShowAliasesRepositoryInterface
	Logger        *slog.Logger
}

func NewHandler(deps *show_aliases.Dependencies) *Handler {
	return &Handler{
		ShowAliasRepo: deps.ShowAliasRepo,
		Logger:        deps.Logger,
	}
}
