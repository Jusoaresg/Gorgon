package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/internal/filter_settings"
	filterSettingsRepository "github.com/jusoaresg/gorgon/internal/filter_settings/repository"
)

type Handler struct {
	FilterSettingsRepo filterSettingsRepository.FilterSettingsRepositoryInterface
	Logger             *slog.Logger
}

func NewHandler(deps *filter_settings.Dependencies) *Handler {
	return &Handler{
		FilterSettingsRepo: deps.FilterSettingsRepo,
		Logger:             deps.Logger,
	}
}
