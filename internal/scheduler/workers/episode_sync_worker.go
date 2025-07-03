package workers

import (
	"log/slog"
	"sync"
	"time"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/episode/model"

	"github.com/jmoiron/sqlx"
)

func StartEpisodeSyncWorker(workerCount int) {
	logger := config.GetLogger().WithGroup("worker").With("name", "episodeSync")
	episodeChan := make(chan model.Episode, 50)
	var wg sync.WaitGroup

	logger.Info("Starting episode sync workers", slog.Int("count", workerCount))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			logger.Info("Worker started")
			processEpisodesWorker(episodeChan)
		}()
	}

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		<-ticker.C
		logger.Info("Checking for wanted episodes")

		episodes := fetchWantedEpisodes()
		if len(episodes) == 0 {
			logger.Info("No episodes found with status 'wanted' or 'missing'")
			continue
		}

		for _, ep := range episodes {
			logger.Info("Queuing episode for processing",
				slog.Int("episodeID", int(ep.ID)),
				slog.Int64("showID", ep.ShowID),
				slog.String("name", ep.Name),
			)
			episodeChan <- ep
		}
	}
}

func fetchWantedEpisodes() []model.Episode {
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
		logger.Error("Failed to build SQL query with sqlx.In", slog.String("error", err.Error()))
		return nil
	}
	query = db.Rebind(query)

	err = db.Select(&episodes, query, args...)
	if err != nil {
		logger.Error("Failed to execute SELECT query", slog.String("error", err.Error()))
		return nil
	}

	return episodes
}
