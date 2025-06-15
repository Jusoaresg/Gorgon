package workers

import (
	"gorgon/internal/db/model"
	"gorgon/pkg/services"
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
	var episodes []model.Episode
	var baseService = services.NewBaseService()
	baseService.DB.Where("tracking = ?", model.StatusSnatched).Limit(100).Find(&episodes)

	return episodes
}
