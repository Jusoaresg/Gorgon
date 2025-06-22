package scheduler

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	episodeRepositoy "github.com/jusoaresg/gorgon/internal/episode/repository"
	epContentRepository "github.com/jusoaresg/gorgon/internal/episode_content/repository"
	"log/slog"
	"os"
	"path/filepath"
)

func VerifyEpisodeWasDeleted() {
	logger := config.GetLogger()
	db := config.GetSQLite()

	configFile, err := config.LoadConfig()
	_ = configFile
	if err != nil {
		return
	}

	episodeRepo := episodeRepositoy.NewEpisodeRepository(db)
	episodes, err := episodeRepo.ListByTracking(model.TrackingDownloaded)
	if err != nil {
		return
	}

	episodeContentRepo := epContentRepository.NewEpisodeContentRepository(db)

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

			fileFolder := filepath.Join("assets", configFile.QBittorrentDownloadFolder, episode_content.Name)
			filePath, _ := filepath.Abs(fileFolder)

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
