package workers

import (
	"gorgon/config"
	prowlarr "gorgon/external/prowlarr/service"
	qbittorrent "gorgon/external/qbittorrent/service"
	"gorgon/internal/db/model"
	"gorgon/internal/scheduler/jobs"
	"log/slog"
)

func processEpisodesWorker(episodesChan <-chan model.Episode) {
	logger := config.GetLogger()
	prowlarrService := prowlarr.NewProwlarrSearchService(logger)
	qbittorrentService, err := qbittorrent.NewQBittorrentService(logger)
	if err != nil {
		logger.Error("Error initializing qbittorrent service", slog.String("error", err.Error()))
		return
	}

	for episode := range episodesChan {
		err := jobs.ProcessSingleEpisode(&episode, *prowlarrService, *qbittorrentService)
		if err != nil {
			logger.Warn("Error processing episode", slog.Int("episodeID", int(episode.ID)), slog.String("error", err.Error()))
		}
	}
}
