package cron

import (
	"gorgon/config"
	"time"
)

func StartDailyUpdate(callback func()) {
	logger := config.GetLogger()
	go func() {
		for {
			logger.Info("Starting to updating shows")
			callback()
			logger.Info("Shows update completed")
			time.Sleep(1 * time.Hour)
		}
	}()
}
