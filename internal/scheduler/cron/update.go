package cron

import (
	"github.com/jusoaresg/gorgon/config"
	"time"
)

func StartDailyUpdate(callback func()) {
	logger := config.GetLogger().WithGroup("scheduler").With("name", "StartDailyUpdate")
	go func() {
		for {
			logger.Info("Starting to updating shows")
			callback()
			logger.Info("Shows update completed")
			time.Sleep(1 * time.Hour)
		}
	}()
}
