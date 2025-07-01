package workers

import (
	"log/slog"
	"time"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	"github.com/jusoaresg/gorgon/internal/scheduler/jobs"
	"github.com/jusoaresg/gorgon/pkg/services"
)

func processSnatchedDownloadsWorker(episodesChan <-chan model.Episode) {
	logger := config.GetLogger().WithGroup("worker").With("name", "processSnatchedDownloadsWorker")

	var qbittorrentService *service.QBittorrentService
	var err error

	for {
		qbittorrentService, err = service.NewQBittorrentService(logger)
		if err != nil {
			logger.Error("Error initializing qbittorrent service", slog.String("error", err.Error()))
			time.Sleep(30 * time.Second)
			continue
		}

		errs := services.CheckAllConnections(qbittorrentService)
		if len(errs) > 0 {
			logger.Error("Failed to connect to one or more services", slog.Any("errors", errs))
			time.Sleep(30 * time.Second)
			continue
		}
		break
	}

	episodeLogger := logger.WithGroup("episode").With("source", "snatched")
	for episode := range episodesChan {
		err := jobs.ProcessSingleSnatchedDownload(&episode, qbittorrentService)
		if err != nil {
			episodeLogger.Error(
				"Error processing snatched episode downloaded",
				slog.Int("episodeID", int(episode.ID)),
				slog.Int64("showID", episode.ShowID),
				slog.String("episodeName", episode.Name),
				slog.String("error", err.Error()))
		}
	}
}
