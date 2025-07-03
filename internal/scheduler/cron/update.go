package cron

import (
	"github.com/jusoaresg/gorgon/config"
	"time"
)

func StartDailyUpdate(callback func()) {
	logger := config.GetLogger().WithGroup("scheduler").With("name", "StartDailyUpdate")
	go func() {
		for {
			logger.Info("starting to updating shows")
			callback()
			logger.Info("shows update completed")
			time.Sleep(24 * time.Hour)
		}
	}()
}
