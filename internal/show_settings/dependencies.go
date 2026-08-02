package show_settings

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showSearchPatternsRepository "github.com/jusoaresg/gorgon/internal/show_search_patterns/repository"
	showSettingsRepository "github.com/jusoaresg/gorgon/internal/show_settings/repository"
)

type Dependencies struct {
	ShowSettingsRepo       showSettingsRepository.ShowSettingsRepositoryInterface
	ShowSearchPatternsRepo showSearchPatternsRepository.ShowSearchPatternsRepositoryInterface
	ShowRepo               showRepository.ShowRepositoryInterface
	DB                     *sqlx.DB
	Logger                 *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	settingsRepo := showSettingsRepository.NewShowSettingsRepository(DB)
	searchPatternsRepo := showSearchPatternsRepository.NewShowSearchPatternsRepository(DB)
	showRepo := showRepository.NewShowRepository(DB)
	return &Dependencies{
		ShowSettingsRepo:       &settingsRepo,
		ShowSearchPatternsRepo: &searchPatternsRepo,
		ShowRepo:               showRepo,
		DB:                     DB,
		Logger:                 logger,
	}
}
