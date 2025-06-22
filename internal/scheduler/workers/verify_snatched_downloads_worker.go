package workers

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/db/model"
	"github.com/jusoaresg/gorgon/internal/db/repository"
	"log/slog"
	"sync"
	"time"
)

func VerifySnatchedDownloadsWorker(workerCount int) {
	episodeChan := make(chan model.Episode, 50)
	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processSnatchedDownloadsWorker(episodeChan)
		}()
	}

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		<-ticker.C

		episodes := fetchSnatchedEpisodes()
		if len(episodes) == 0 {
			continue
		}

		for _, ep := range episodes {
			episodeChan <- ep
		}

	}
}

func fetchSnatchedEpisodes() []model.Episode {
	logger := config.GetLogger()
	episodeRepo := repository.NewEpisodeRepository(config.GetSQLite())

	//TODO: Put an limit of 100 here
	episodes, err := episodeRepo.ListByTracking(model.TrackingSnatched)
	if err != nil {
		logger.Error(
			"Error fetching snatched episodes",
			slog.Any("Episodes", episodes),
			slog.String("worker", "fetchSnatchedEpisodes"),
			slog.String("error", err.Error()))
		return nil
	}

	return episodes
}
