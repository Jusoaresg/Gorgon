package workers

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/config"
	prowlarr "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrent "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	"github.com/jusoaresg/gorgon/internal/scheduler/jobs"
	"github.com/jusoaresg/gorgon/pkg/services"
)

func processEpisodesWorker(episodesChan <-chan model.Episode) {
	logger := config.GetLogger().WithGroup("worker").With("name", "processEpisodeWorker")
	prowlarrService, err := prowlarr.NewProwlarrSearchService(logger)
	if err != nil {
		logger.Error("Error initializing prowlarr service", slog.String("error", err.Error()))
		return
	}
	qbittorrentService, err := qbittorrent.NewQBittorrentService(logger)
	if err != nil {
		logger.Error("Error initializing qbittorrent service", slog.String("error", err.Error()))
		return
	}

	errs := services.CheckAllConnections(prowlarrService, qbittorrentService)
	if len(errs) > 0 {
		logger.Error("Failed to connect to one or more services", slog.Any("errors", err))
		return
	}

	episodeLogger := logger.WithGroup("episode")
	for episode := range episodesChan {
		err := jobs.ProcessSingleEpisode(&episode, *prowlarrService, *qbittorrentService)
		if err != nil {
			episodeLogger.Error(
				"Error processing single episode",
				slog.Int("episodeID", int(episode.ID)),
				slog.Int64("showID", episode.ShowID),
				slog.String("episodeName", episode.Name),
				slog.String("error", err.Error()))
		}
	}
}
