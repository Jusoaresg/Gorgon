package service

import (
	"log/slog"
	"sync"

	"github.com/jusoaresg/gorgon/config"
	episodeEvents "github.com/jusoaresg/gorgon/internal/episode/events"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	"github.com/jusoaresg/gorgon/internal/episode/repository"
	"github.com/jusoaresg/gorgon/pkg/concurrency"
)

var DbWriteMutex sync.Mutex

func SnatchEpisode(episode *model.Episode, torrentHash string) error {
	logger := config.GetLogger()
	safeDB := config.GetSafeDB()

	episodeRepo := repository.NewEpisodeRepository(safeDB.Db)

	episode.Tracking = model.TrackingSnatched
	episode.TorrentHash = torrentHash

	err := concurrency.WithWriteLock(safeDB.Write, func() error {
		return episodeRepo.Update(*episode)
	})
	if err != nil {
		logger.Error(
			"failed to update episode status to snatched",
			slog.Int64("episode_id", episode.ID),
			slog.Int64("show_id", episode.ShowID),
			slog.String("torrent_hash", torrentHash),
		)
		return err
	}

	episodeEvents.EmitEpisodeTrackingUpdatedEvent(episode.ID, model.TrackingSnatched)

	return nil
}
