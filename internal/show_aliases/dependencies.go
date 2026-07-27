package show_aliases

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
)

type Dependencies struct {
	ShowAliasRepo showAliasRepository.ShowAliasesRepositoryInterface
	DB            *sqlx.DB
	Logger        *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	repo := showAliasRepository.NewShowAliasesRepository(DB)
	return &Dependencies{
		ShowAliasRepo: &repo,
		DB:            DB,
		Logger:        logger,
	}
}
