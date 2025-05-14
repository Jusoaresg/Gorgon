package cron

import (
	"gorgon/config"
	"time"
)

func StartVerifySnatched(callback func()) {
	logger := config.GetLogger()
	go func() {
		for {
			logger.Info("Starting to verifying if any snatched has been downloaded")
			callback()
			logger.Info("Snatched verification completed")
			time.Sleep(30 * time.Second)
		}
	}()
}
