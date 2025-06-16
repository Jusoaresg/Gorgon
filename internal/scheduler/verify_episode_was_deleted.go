package scheduler

import (
	"fmt"
	"gorgon/config"
	"gorgon/internal/db/model"
	"gorgon/internal/db/repository"
	"log/slog"
	"os"
)

func VerifyEpisodeWasDeleted() {
	logger := config.GetLogger()

	configFile, err := config.LoadConfig()
	_ = configFile
	if err != nil {
		return
	}

	episodeRepo := repository.NewEpisodeRepository()
	episodes, err := episodeRepo.ListByTracking(model.TrackingDownloaded)
	if err != nil {
		return
	}

	episodeContentRepo := repository.NewEpisodeContentRepository()

	for _, episode := range episodes {
		contents, err := episodeContentRepo.ListByEpisodeId(episode.ID)
		if err != nil {
			continue
		}

		if contents == nil {
			continue
		}

		logger.Debug("Downloaded episode found", slog.Int("Episode Number", episode.Number), slog.String("Episode Name", episode.Name))
		for _, episode_content := range contents {

			filePath := fmt.Sprintf("/home/juliano/Videos/downloads/%s", episode_content.Name)

			_, err := os.Stat(filePath)
			if os.IsNotExist(err) {
				episode.SetNotInstalled()

				err := episodeRepo.Update(episode)
				if err != nil {
					continue
				}

				episodeContentRepo.DeleteById(episode_content.ID)

				logger.Info("Episode file not found, setting tracking to skipped", slog.Int("Episode Number", episode.Number), slog.String("Episode Name", episode.Name))
			}
		}
	}
}
