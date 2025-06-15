package jobs

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/qbittorrent/schema"
	"gorgon/external/qbittorrent/service"
	"gorgon/internal/db/model"
	"gorgon/pkg/handler"
	"gorgon/pkg/services"
)

type EpisodeUpdatedWebsocketSchema struct {
	Type      string `json:"type"`
	EpisodeId int    `json:"episodeId"`
	Tracking  string `json:"tracking"`
}

func ProcessSingleSnatchedDownload(ep *model.Episode, qbittorrentService service.QBittorrentService) error {
	logger := config.GetLogger()

	var baseService = services.NewBaseService()
	var torrentResponse []schema.CheckTorrentResponse
	qbittorrentService.CheckTorrentsWithHash("completed", ep.TorrentHash, &torrentResponse)
	if len(torrentResponse) == 0 {
		logger.Info(fmt.Sprintf("Episode S%02d E%02d %s not found between the completed torrents.", ep.Season, ep.Number, ep.Name))
		return nil
	}
	torrent := torrentResponse[0]

	if torrent.Hash == ep.TorrentHash {
		logger.Info(fmt.Sprintf("Episode S%02d E%02d %s found - Torrent: %s - State: %s - Progress: %.2f%%", ep.Season, ep.Number, ep.Name, torrent.Name, torrent.State, torrent.Progress*100))

		updated_episode := ep
		updated_episode.Tracking = model.Tracking.Downloaded()
		updated_episode.FilePath = &torrent.SavePath
		qbittorrentService.CheckContent(torrent.Hash, &updated_episode.Content)

		baseService.UpdateByID(int(ep.ID), &updated_episode)

		var msg EpisodeUpdatedWebsocketSchema
		msg.Type = "EpisodeTrackingUpdate"
		msg.EpisodeId = int(ep.ID)
		msg.Tracking = string(model.Tracking.Downloaded())
		handler.SendWebSocketMessage(msg)
		return nil
	}
	logger.Info(fmt.Sprintf("Episode S%02d E%02d %s not found between the progress torrents.", ep.Season, ep.Number, ep.Name))
	return nil
}
