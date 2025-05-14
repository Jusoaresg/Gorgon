package cron

import (
	"gorgon/config"
	"time"
)

func StartSearchNewEpisodes(callback func()) {
	logger := config.GetLogger()
	go func() {
		for {
			logger.Info("Starting to searching for availble episodes")
			callback()
			logger.Info("Episode search completed")
			time.Sleep(1 * time.Minute)
		}
	}()
}
