package app

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/downloads"
	"github.com/jusoaresg/gorgon/internal/episode"
	"github.com/jusoaresg/gorgon/internal/indexer"
	appConfig "github.com/jusoaresg/gorgon/internal/config"
	"github.com/jusoaresg/gorgon/internal/season"
	"github.com/jusoaresg/gorgon/internal/show"
	"github.com/jusoaresg/gorgon/internal/show_aliases"
)

type Dependencies struct {
	Show        *show.Dependencies
	Episode     *episode.Dependencies
	Season      *season.Dependencies
	Indexer     *indexer.Dependencies
	Downloads   *downloads.Dependencies
	AppConfig   *appConfig.Dependencies
	ShowAliases *show_aliases.Dependencies
	Logger      *slog.Logger
	DB          *sqlx.DB
}

func NewDependencies() *Dependencies {
	logger := config.GetLogger()
	db := config.GetSQLite()

	return &Dependencies{
		Show:        show.NewDependencies(db, logger),
		Episode:     episode.NewDependencies(db, logger),
		Season:      season.NewDependencies(db, logger),
		Indexer:     indexer.NewDependencies(db, logger),
		Downloads:   downloads.NewDependencies(db, logger),
		AppConfig:   appConfig.NewDependencies(logger),
		ShowAliases: show_aliases.NewDependencies(db, logger),
		Logger:      logger,
		DB:          db,
	}
}
