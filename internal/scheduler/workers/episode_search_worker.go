package workers

import (
	"log/slog"
	"sync"
	"time"

	"github.com/jusoaresg/gorgon/config"
	prowlarr "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrent "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	"github.com/jusoaresg/gorgon/internal/episode/repository"
)

func StartEpisodeSearchWorker(workerCount int, prowlarrService *prowlarr.ProwlarrSearchService, qbittorrentService *qbittorrent.QBittorrentService) {
	logger := config.GetLogger().WithGroup("worker").With("name", "episodeSync")

	episodeChan := make(chan model.Episode, 50)
	var wg sync.WaitGroup

	logger.Info("Starting episode sync workers", slog.Int("count", workerCount))

	for i := range workerCount {
		workerID := i
		wg.Add(1)
		wg.Go(func() {
			defer wg.Done()
			logger.Info("worker started", slog.Int("worker_id", workerID))
			processEpisodesWorker(episodeChan, prowlarrService, qbittorrentService)
		})
	}

	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for {
		<-ticker.C
		logger.Info("Checking for wanted episodes")

		episodes := fetchEpisodeSearchWantedEpisodes()
		if len(episodes) == 0 {
			logger.Info("No episodes found with status 'wanted' or 'missing'")
			continue
		}

		for _, ep := range episodes {
			logger.Info("Queuing episode for processing",
				slog.Int("episode_id", int(ep.ID)),
				slog.Int64("show_id", ep.ShowID),
				slog.String("name", ep.Name),
			)
			episodeChan <- ep
		}
	}
}

func fetchEpisodeSearchWantedEpisodes() []model.Episode {
	logger := config.GetLogger().WithGroup("worker").With("name", "episodeSync")
	db := config.GetSQLite()

	episodeRepo := repository.NewEpisodeRepository(db)

	episodes, err := episodeRepo.ListReleasedByTracking("wanted", "missing")
	if err != nil {
		logger.Error("Failed to list releases by tracking", slog.String("error", err.Error()))
		return nil
	}

	return episodes
}
