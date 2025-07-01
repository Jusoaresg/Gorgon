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
	logger := config.GetLogger().WithGroup("scheduler").With("name", "VerifyEpisodeWasDeleted")
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

		for _, episode_content := range contents {

			fileFolder := filepath.Join(configFile.QBittorrentDownloadFolder, episode_content.Name)
			filePath, _ := filepath.Abs(fileFolder)

			_, err := os.Stat(filePath)
			if os.IsNotExist(err) {
				episode.SetNotInstalled()

				err := episodeRepo.Update(episode)
				if err != nil {
					continue
				}

				episodeContentRepo.DeleteById(episode_content.ID)

				logger.Info(
					"Episode file not found, setting tracking to skipped",
					slog.Int64("showID", episode.ShowID),
					slog.Int("episodeNumber", episode.Number),
					slog.String("episodeName", episode.Name),
				)
			}
		}
	}
}
