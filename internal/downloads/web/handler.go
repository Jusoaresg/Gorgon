package web

import (
	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/internal/downloads"
	"github.com/jusoaresg/gorgon/internal/downloads/service"
)

type Handler struct {
	Service *service.DownloadsService
	DB      *sqlx.DB
}

func NewHandler(deps *downloads.Dependencies) *Handler {
	return &Handler{
		Service: deps.Service,
		DB:      deps.DB,
	}
}
