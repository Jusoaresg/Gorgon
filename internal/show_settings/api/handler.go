package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/internal/show_settings"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showSearchPatternsRepository "github.com/jusoaresg/gorgon/internal/show_search_patterns/repository"
	showSettingsRepository "github.com/jusoaresg/gorgon/internal/show_settings/repository"
)

type Handler struct {
	ShowSettingsRepo       showSettingsRepository.ShowSettingsRepositoryInterface
	ShowSearchPatternsRepo showSearchPatternsRepository.ShowSearchPatternsRepositoryInterface
	ShowRepo               showRepository.ShowRepositoryInterface
	Logger                 *slog.Logger
}

func NewHandler(deps *show_settings.Dependencies) *Handler {
	return &Handler{
		ShowSettingsRepo:       deps.ShowSettingsRepo,
		ShowSearchPatternsRepo: deps.ShowSearchPatternsRepo,
		ShowRepo:               deps.ShowRepo,
		Logger:                 deps.Logger,
	}
}
