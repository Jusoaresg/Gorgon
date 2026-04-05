package workers

import (
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/external/prowlarr/service"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	"github.com/jusoaresg/gorgon/internal/episode/repository"
)

func StartRssEpisodeFetcherWorker(workerCount int, prowlarrService *service.ProwlarrSearchService) {
	logger := config.GetLogger().WithGroup("worker").With("name", "StartRssFeedWorker")
	rssProcessor := NewRssReleaseProcessor(config.GetSQLite())

	episodeChan := make(chan model.Episode)
	var responses []schema.SearchResponse

	var wg sync.WaitGroup

	logger.Info(
		"Starting rss feed worker",
		slog.Int("worker_count", workerCount),
	)

	for i := range workerCount {
		workerID := i
		wg.Add(1)
		wg.Go(func() {
			defer wg.Done()
			logger.Info("Worker started", slog.Int("worker_id", workerID))

			for ep := range episodeChan {
				err := rssProcessor.RssProcessRelease(ep, responses)
				if err != nil {
					logger.Error(
						"failed to process episode",
						slog.String("error", err.Error()),
						slog.Int64("show_id", ep.ShowID),
						slog.Int64("episode_id", ep.ID),
					)
					continue
				}
			}
		})
	}

	ticker := time.NewTicker(time.Second * 45)
	defer ticker.Stop()

	for {
		<-ticker.C
		logger.Info("checking for wanted episodes")

		episodes := fetchRssFeedWantedEpisodes()
		if len(episodes) == 0 {
			logger.Info("no episodes found with status 'wanted' or 'missing'")
			continue
		}

		//TODO: Required words here
		query := strings.Join([]string{"multi subs"}, " ")
		prowlarrService.Search(&schema.SearchRequest{Query: query}, &responses)

		for _, ep := range episodes {
			logger.Info(
				"queuing episode for processing",
				slog.Int64("episode_id", ep.ID),
				slog.Int64("show_id", ep.ShowID),
				slog.String("name", ep.Name),
			)
			episodeChan <- ep
		}
	}
}

func fetchRssFeedWantedEpisodes() []model.Episode {
	logger := config.GetLogger().WithGroup("worker").With("name", "fetchWantedEps")
	db := config.GetSQLite()

	episodeRepositoy := repository.NewEpisodeRepository(db)

	episodes, err := episodeRepositoy.ListByTracking("wanted", "missing")
	if err != nil {
		logger.Error("Failed to list episodes by tracking", slog.String("error", err.Error()))
	}
	return episodes
}
