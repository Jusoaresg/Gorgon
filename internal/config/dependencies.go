package config

import (
	"log/slog"
)

type Dependencies struct {
	Logger *slog.Logger
}

func NewDependencies(logger *slog.Logger) *Dependencies {
	return &Dependencies{
		Logger: logger,
	}
}
