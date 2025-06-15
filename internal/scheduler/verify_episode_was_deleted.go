package scheduler

import (
	"fmt"
	"gorgon/config"
	"gorgon/internal/db/model"
	"gorgon/pkg/services"
	"log/slog"
	"os"
)

func VerifyEpisodeWasDeleted() {
	logger := config.GetLogger()

	baseService := services.NewBaseService()
	configFile, err := config.LoadConfig()
	_ = configFile
	if err != nil {
		return
	}

	var episodes []model.Episode
	baseService.DB.Preload("Content").Where("tracking = ?", "downloaded").Find(&episodes)

	for _, episode := range episodes {
		if episode.Tracking != "downloaded" {
			logger.Debug("No downloaded episode found", slog.Int("Episode Number", episode.Number), slog.String("Episode Name", episode.Name))
			continue
		}

		if episode.Content == nil {
			continue
		}

		if episode.Tracking == "downloaded" && episode.Content != nil {

			logger.Debug("Downloaded episode found", slog.Int("Episode Number", episode.Number), slog.String("Episode Name", episode.Name))
			for _, episode_content := range episode.Content {

				filePath := fmt.Sprintf("/home/juliano/Videos/downloads/%s", episode_content.Name)

				_, err := os.Stat(filePath)
				if os.IsNotExist(err) {
					episode.SetNotInstalled()

					baseService.ListWithIdentification(&episode_content, "episode_id", string(episode.ID))

					baseService.UpdateByIDWithSelect(int(episode.ID), &episode, "FilePath", "TorrentHash", "Tracking")
					baseService.Delete(int(episode_content.ID), &episode_content)

					logger.Info("Episode file not found, setting tracking to skipped", slog.Int("Episode Number", episode.Number), slog.String("Episode Name", episode.Name))
				}
			}
		}
	}
}
