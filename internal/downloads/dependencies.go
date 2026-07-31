package downloads

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/internal/downloads/service"
)

type Dependencies struct {
	Service *service.DownloadsService
	DB      *sqlx.DB
	Logger  *slog.Logger
}

func NewDependencies(DB *sqlx.DB, logger *slog.Logger) *Dependencies {
	return &Dependencies{
		Service: service.NewDownloadsService(DB, logger),
		DB:      DB,
		Logger:  logger,
	}
}
