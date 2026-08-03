package api

import (
	"log/slog"

	prowlarrService "github.com/jusoaresg/gorgon/external/prowlarr/service"
	"github.com/jusoaresg/gorgon/internal/indexer"
	"github.com/jusoaresg/gorgon/internal/indexer/repository"
)

type Handler struct {
	IndexerRepo        *repository.IndexerRepository
	ProwlarrIndexerSvc *prowlarrService.ProwlarrIndexerService
	Logger             *slog.Logger
}

func NewHandler(deps *indexer.Dependencies) *Handler {
	return &Handler{
		IndexerRepo:        deps.IndexerRepo,
		ProwlarrIndexerSvc: deps.ProwlarrIndexerSvc,
		Logger:             deps.Logger,
	}
}
