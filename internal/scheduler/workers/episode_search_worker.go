package workers

import (
	"log/slog"
	"sync"
	"time"

	"github.com/jusoaresg/gorgon/config"
	prowlarr "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrent "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/episode/model"

	"github.com/jmoiron/sqlx"
)

func StartEpisodeSearchWorker(workerCount int, prowlarrService *prowlarr.ProwlarrSearchService, qbittorrentService *qbittorrent.QBittorrentService) {
	logger := config.GetLogger().WithGroup("worker").With("name", "episodeSync")
	episodeChan := make(chan model.Episode, 50)
	var wg sync.WaitGroup

	logger.Info("starting episode sync workers", slog.Int("count", workerCount))

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
		logger.Info("checking for wanted episodes")

		episodes := fetchEpisodeSearchWantedEpisodes()
		if len(episodes) == 0 {
			logger.Info("no episodes found with status 'wanted' or 'missing'")
			continue
		}

		for _, ep := range episodes {
			logger.Info("queuing episode for processing",
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
	var episodes []model.Episode
	db := config.GetSQLite()

	query, args, err := sqlx.In(
		`SELECT * FROM episodes WHERE tracking IN (?) 
			AND airstamp <= ? 
		LIMIT 20`,
		[]string{"wanted", "missing"},
		time.Now().UTC().Unix(),
	)
	if err != nil {
		logger.Error("failed to build SQL query with sqlx.In", slog.String("error", err.Error()))
		return nil
	}
	query = db.Rebind(query)

	err = db.Select(&episodes, query, args...)
	if err != nil {
		logger.Error("failed to execute SELECT query", slog.String("error", err.Error()))
		return nil
	}

	return episodes
}
