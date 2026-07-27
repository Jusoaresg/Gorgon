package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/internal/config"
)

type Handler struct {
	Logger *slog.Logger
}

func NewHandler(deps *config.Dependencies) *Handler {
	return &Handler{
		Logger: deps.Logger,
	}
}
