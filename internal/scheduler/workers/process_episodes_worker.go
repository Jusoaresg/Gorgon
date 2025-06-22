package workers

import (
	"github.com/jusoaresg/gorgon/config"
	prowlarr "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrent "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/db/model"
	"github.com/jusoaresg/gorgon/internal/scheduler/jobs"
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
			logger.Error(
				"Error processing episode",
				slog.Int("episodeID", int(episode.ID)),
				slog.Int64("showID", episode.ShowID),
				slog.String("episodeName", episode.Name),
				slog.String("worker", "processEpisodeWorker"),
				slog.String("error", err.Error()))
		}
	}
}
