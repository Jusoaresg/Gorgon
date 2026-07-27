package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/internal/season"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
)

type Handler struct {
	SeasonRepo seasonRepository.SeasonRepositoryInterface
	Logger     *slog.Logger
}

func NewHandler(deps *season.Dependencies) *Handler {
	return &Handler{
		SeasonRepo: deps.SeasonRepo,
		Logger:     deps.Logger,
	}
}
