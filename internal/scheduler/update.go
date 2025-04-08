package scheduler

import (
	"gorgon/config"
	"os"
	"time"
)

const updateFile = "assets/last_update.txt"

func StartDailyUpdate(callback func()) {
	logger := config.GetLogger()
	go func() {
		for {
			// now := time.Now()
			logger.Debug("Checking if shows has already updated today")
			if !alreadyUpdatedToday() {
				logger.Info("Starting to updating shows")
				callback()
				saveTodayAsUpdated()
				logger.Info("Shows update completed")
			}
			time.Sleep(1 * time.Hour)
		}
	}()
}

func alreadyUpdatedToday() bool {
	data, err := os.ReadFile(updateFile)
	if err != nil {
		return false
	}
	lastUpdate := string(data)
	today := time.Now().Format("2006-01-02")
	return lastUpdate == today
}

func saveTodayAsUpdated() {
	today := time.Now().Format("2006-01-02")
	_ = os.WriteFile(updateFile, []byte(today), 0644)
}
