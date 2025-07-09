package scheduler

import "github.com/jusoaresg/gorgon/internal/scheduler/workers"

func Start() {
	go workers.StartRssFeedWorker(5)
	go workers.StartEpisodeSearchWorker(2)

	go workers.VerifySnatchedDownloadsWorker(5)
}
