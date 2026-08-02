package filter_settings

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	filterSettingsRepository "github.com/jusoaresg/gorgon/internal/filter_settings/repository"
)

type Dependencies struct {
	FilterSettingsRepo filterSettingsRepository.FilterSettingsRepositoryInterface
	DB                 *sqlx.DB
	Logger             *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	repo := filterSettingsRepository.NewFilterSettingsRepository(DB)
	return &Dependencies{
		FilterSettingsRepo: &repo,
		DB:                 DB,
		Logger:             logger,
	}
}
