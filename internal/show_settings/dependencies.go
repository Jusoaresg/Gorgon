package show_settings

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showSettingsRepository "github.com/jusoaresg/gorgon/internal/show_settings/repository"
)

type Dependencies struct {
	ShowSettingsRepo showSettingsRepository.ShowSettingsRepositoryInterface
	ShowRepo         showRepository.ShowRepositoryInterface
	DB               *sqlx.DB
	Logger           *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	settingsRepo := showSettingsRepository.NewShowSettingsRepository(DB)
	showRepo := showRepository.NewShowRepository(DB)
	return &Dependencies{
		ShowSettingsRepo: &settingsRepo,
		ShowRepo:         showRepo,
		DB:               DB,
		Logger:           logger,
	}
}
