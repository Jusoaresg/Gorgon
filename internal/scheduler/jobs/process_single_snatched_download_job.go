package jobs

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/qbittorrent/schema"
	"gorgon/external/qbittorrent/service"
	"gorgon/internal/db/events/episode"
	"gorgon/internal/db/model"
	"gorgon/internal/db/repository"
	"gorgon/internal/paths"
	"gorgon/utils"
	"log/slog"
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
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

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

		contents, err := qbittorrentService.GetContent(torrent.Hash)
		if err != nil {
			return err
		}

		showRepo := repository.NewShowRepository(config.GetSQLite())
		show, err := showRepo.GetById(ep.ShowID)

		tx, err := config.GetSQLite().Beginx()
		if err != nil {
			return err
		}
		episodeRepo.UpdateTx(tx, *ep)
		for _, content := range contents {
			content.FilePath = torrent.SavePath
			content.EpisodeId = ep.ID
			if err := episodeContentRepo.CreateTx(tx, content); err != nil {
				return err
			}

			symlinkPath, err := utils.SymlinkPathForEpisode(cfg.ShowsFolder, show.Name, *ep, content)
			if err != nil {
				return err
			}
			episodeDownloadFolder, err := paths.GetEpisodeDownloadFile(content.Name)
			if err != nil {
				return err
			}

			if err := utils.CreateSymlink(episodeDownloadFolder, symlinkPath); err != nil {
				logger.Error("Failed to create symlink", slog.String("from", episodeDownloadFolder), slog.String("to", symlinkPath), slog.String("error", err.Error()))
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		episode.EmitEpisodeTrackingUpdatedEvent(ep.ID, model.TrackingDownloaded)

		return nil
	}
	logger.Info(fmt.Sprintf("Episode S%02d E%02d %s not found between the progress torrents.", ep.Season, ep.Number, ep.Name))
	return nil
}
