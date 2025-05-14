package scheduler

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/qbittorrent/schema"
	"gorgon/external/qbittorrent/service"
	"gorgon/internal/db/model"
	"gorgon/pkg/services"
)

func VerifySnatchedDownload() {
	logger := config.GetLogger()

	qbittorrentService, err := service.NewQBittorrentService(logger)
	if err != nil {
		//TODO: Log error message
		return
	}
	baseService := services.NewBaseService()

	var episodes []model.Episode
	baseService.ListWithIdentification(&episodes, "tracking", string(model.Tracking.Snatched()))
	if len(episodes) <= 0 {
		//TODO: Log message
		return
	}

	var torrentResponse []schema.CheckTorrentResponse
	qbittorrentService.CheckTorrents("completed", &torrentResponse)

	if len(torrentResponse) <= 0 {
		return
	}

	for _, episode := range episodes {
		//TODO: An option to remove the torrent after download would be good
		episodeMatched := false

		for _, torrent := range torrentResponse {
			if episode.TorrentHash == torrent.Hash {
				episodeMatched = true
				logger.Info(fmt.Sprintf("Episode S%02d E%02d %s found - Torrent: %s - State: %s - Progress: %.2f%%", episode.Season, episode.Number, episode.Name, torrent.Name, torrent.State, torrent.Progress*100))

				updated_episode := episode
				updated_episode.Tracking = model.Tracking.Downloaded()
				updated_episode.FilePath = &torrent.SavePath
				baseService.UpdateByID(int(episode.ID), &updated_episode)
				break
			}
		}
		if !episodeMatched {
			logger.Info(fmt.Sprintf("Episode S%02d E%02d %s not found between the progress torrents.", episode.Season, episode.Number, episode.Name))
		}
	}
}
