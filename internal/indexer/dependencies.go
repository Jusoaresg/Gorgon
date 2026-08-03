package indexer

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	prowlarrService "github.com/jusoaresg/gorgon/external/prowlarr/service"
	"github.com/jusoaresg/gorgon/internal/indexer/repository"
)

type Dependencies struct {
	IndexerRepo          *repository.IndexerRepository
	ProwlarrIndexerSvc   *prowlarrService.ProwlarrIndexerService
	Logger               *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	return &Dependencies{
		IndexerRepo:        repository.NewIndexerRepository(DB),
		ProwlarrIndexerSvc: prowlarrService.NewProwlarrIndexerService(logger),
		Logger:             logger,
	}
}
