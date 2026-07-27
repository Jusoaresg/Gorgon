package season

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
)

type Dependencies struct {
	SeasonRepo seasonRepository.SeasonRepositoryInterface
	DB         *sqlx.DB
	Logger     *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	return &Dependencies{
		SeasonRepo: seasonRepository.NewSeasonRepository(DB),
		DB:         DB,
		Logger:     logger,
	}
}
