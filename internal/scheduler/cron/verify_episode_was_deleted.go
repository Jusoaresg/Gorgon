package cron

import (
	"github.com/jusoaresg/gorgon/config"
	"time"
)

func StartVerifyEpisodeWasDeleted(callback func()) {
	logger := config.GetLogger().WithGroup("cron").With("name", "StartVerifyEpisodeWasDeleted")
	go func() {
		for {
			logger.Info("Starting to verifying if any episode has been deleted")
			callback()
			logger.Info("Verification of deleted episodes completed")
			time.Sleep(30 * time.Second)
		}
	}()
}
