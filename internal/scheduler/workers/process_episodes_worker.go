package workers

import (
	"log/slog"
	"time"

	"github.com/jusoaresg/gorgon/config"
	prowlarr "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrent "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	"github.com/jusoaresg/gorgon/internal/scheduler/jobs"
	"github.com/jusoaresg/gorgon/pkg/services"
)

func processEpisodesWorker(episodesChan <-chan model.Episode) {
	logger := config.GetLogger().WithGroup("worker").With("name", "processEpisodeWorker")

	var prowlarrService *prowlarr.ProwlarrSearchService
	var qbittorrentService *qbittorrent.QBittorrentService
	var err error

	for {
		prowlarrService, err = prowlarr.NewProwlarrSearchService(logger)
		if err != nil {
			logger.Error("Error initializing prowlarr service", slog.String("error", err.Error()))
			time.Sleep(30 * time.Second)
			continue
		}

		qbittorrentService, err = qbittorrent.NewQBittorrentService(logger)
		if err != nil {
			logger.Error("Error initializing qbittorrent service", slog.String("error", err.Error()))
			time.Sleep(30 * time.Second)
			continue
		}

		errs := services.CheckAllConnections(prowlarrService, qbittorrentService)
		if len(errs) > 0 {
			logger.Error("Failed to connect to one or more services", slog.Any("errors", errs))
			time.Sleep(30 * time.Second)
			continue
		}
		break
	}

	episodeLogger := logger.WithGroup("episode")
	for episode := range episodesChan {
		err := jobs.ProcessSingleEpisode(&episode, prowlarrService, qbittorrentService)
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
