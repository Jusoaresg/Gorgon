package workers

import (
	"gorgon/config"
	"gorgon/internal/db/model"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
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
	db := config.GetSQLite()

	query, args, err := sqlx.In("SELECT * FROM episodes WHERE tracking IN (?) LIMIT 20", []string{"wanted", "missing"})
	if err != nil {
		return nil
	}
	query = db.Rebind(query)

	err = db.Select(&episodes, query, args...)
	if err != nil {
		return nil
	}

	return episodes
}
