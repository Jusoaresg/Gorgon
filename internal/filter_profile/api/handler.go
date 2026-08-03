package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/internal/filter_profile"
	filterProfileRepository "github.com/jusoaresg/gorgon/internal/filter_profile/repository"
)

type Handler struct {
	FilterProfileRepo filterProfileRepository.FilterProfileRepositoryInterface
	Logger            *slog.Logger
}

func NewHandler(deps *filter_profile.Dependencies) *Handler {
	return &Handler{
		FilterProfileRepo: deps.FilterProfileRepo,
		Logger:            deps.Logger,
	}
}
