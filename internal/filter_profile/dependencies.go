package filter_profile

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	filterProfileRepository "github.com/jusoaresg/gorgon/internal/filter_profile/repository"
)

type Dependencies struct {
	FilterProfileRepo filterProfileRepository.FilterProfileRepositoryInterface
	DB                *sqlx.DB
	Logger            *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	repo := filterProfileRepository.NewFilterProfileRepository(DB)
	return &Dependencies{
		FilterProfileRepo: &repo,
		DB:                DB,
		Logger:            logger,
	}
}
