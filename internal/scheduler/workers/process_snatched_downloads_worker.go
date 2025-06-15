package workers

import (
	"gorgon/config"
	"gorgon/external/qbittorrent/service"
	"gorgon/internal/db/model"
	"gorgon/internal/scheduler/jobs"
	"log/slog"
)

func processSnatchedDownloadsWorker(episodesChan <-chan model.Episode) {
	logger := config.GetLogger()
	qbittorrentService, err := service.NewQBittorrentService(logger)
	if err != nil {
		logger.Error("Error initializing qbittorrent service", slog.String("error", err.Error()))
		return
	}

	for episode := range episodesChan {
		err := jobs.ProcessSingleSnatchedDownload(&episode, *qbittorrentService)
		if err != nil {
			logger.Warn("Error processing episode", slog.Int("episodeID", int(episode.ID)), slog.String("error", err.Error()))
		}
	}
}
