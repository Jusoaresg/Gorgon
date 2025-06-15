package workers

import (
	"gorgon/internal/db/model"
	"gorgon/pkg/services"
	"sync"
	"time"
)

func StartEpisodeSyncWorker(workerCount int) {
	episodeChan := make(chan model.Episode, 50)
	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processEpisodesWorker(episodeChan)
		}()
	}

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		<-ticker.C

		episodes := fetchWantedEpisodes()
		if len(episodes) == 0 {
			continue
		}

		for _, ep := range episodes {
			episodeChan <- ep
		}
	}
}

func fetchWantedEpisodes() []model.Episode {
	var episodes []model.Episode
	baseService := services.NewBaseService()
	baseService.DB.Where("tracking IN ?", []string{"wanted", "missing"}).Limit(20).Find(&episodes)

	return episodes
}
