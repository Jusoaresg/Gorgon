package jobs

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/qbittorrent/schema"
	"gorgon/external/qbittorrent/service"
	"gorgon/internal/db/model"
	"gorgon/internal/db/repository"
	"gorgon/pkg/handler"
)

type EpisodeUpdatedWebsocketSchema struct {
	Type      string `json:"type"`
	EpisodeId int64  `json:"episodeId"`
	Tracking  string `json:"tracking"`
}

func ProcessSingleSnatchedDownload(ep *model.Episode, qbittorrentService service.QBittorrentService) error {
	logger := config.GetLogger()
	episodeRepo := repository.NewEpisodeRepository(config.GetSQLite())
	episodeContentRepo := repository.NewEpisodeContentRepository()

	var torrentResponse []schema.CheckTorrentResponse
	qbittorrentService.CheckTorrentsWithHash("completed", ep.TorrentHash, &torrentResponse)
	if len(torrentResponse) == 0 {
		logger.Info(fmt.Sprintf("Episode S%02d E%02d %s not found between the completed torrents.", ep.Season, ep.Number, ep.Name))
		return nil
	}
	torrent := torrentResponse[0]

	if torrent.Hash == ep.TorrentHash {
		logger.Info(fmt.Sprintf("Episode S%02d E%02d %s found - Torrent: %s", ep.Season, ep.Number, ep.Name, torrent.Name))

		ep.Tracking = model.TrackingDownloaded
		ep.FilePath = torrent.SavePath

		contents, err := qbittorrentService.GetContent(torrent.Hash)
		if err != nil {
			return err
		}

		tx, err := config.GetSQLite().Beginx()
		if err != nil {
			return err
		}
		episodeRepo.UpdateTx(tx, *ep)
		for _, content := range contents {
			content.EpisodeId = ep.ID
			if err := episodeContentRepo.CreateTx(tx, content); err != nil {
				return err
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		var msg EpisodeUpdatedWebsocketSchema
		msg.Type = "EpisodeTrackingUpdate"
		msg.EpisodeId = ep.ID
		msg.Tracking = model.TrackingDownloaded
		handler.SendWebSocketMessage(msg)
		return nil
	}
	logger.Info(fmt.Sprintf("Episode S%02d E%02d %s not found between the progress torrents.", ep.Season, ep.Number, ep.Name))
	return nil
}
