package jobs

import (
	"github.com/jusoaresg/gorgon/config"
	prowlarr "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrent "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	episodeService "github.com/jusoaresg/gorgon/internal/episode/service"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
)

// var cooldownCache sync.Map
//
// const cooldownDuration = 5 * time.Minute
//
// func init() {
// 	go func() {
// 		ticker := time.NewTicker(30 * time.Minute)
// 		for range ticker.C {
// 			now := time.Now()
// 			cooldownCache.Range(func(key, value any) bool {
// 				if t, ok := value.(time.Time); ok && now.Sub(t) > cooldownDuration*2 {
// 					cooldownCache.Delete(key)
// 				}
// 				return true
// 			})
// 		}
// 	}()
// }

func ProcessBackfillSingleEpisode(ep *model.Episode, prowlarrService *prowlarr.ProwlarrSearchService, qbittorrentService *qbittorrent.QBittorrentService) error {
	logger := config.GetLogger().WithGroup("jobs").With("name", "ProcessSingleEpisode")
	db := config.GetSQLite()

	episodeRepo := episodeRepository.NewEpisodeRepository(db)
	showRepo := showRepository.NewShowRepository(db)

	episodeSearchService := episodeService.NewEpisodeSearchService(
		db,
		logger,
		&episodeService.EpisodeSearcher{},
		&episodeService.EpisodeDownloader{},
		episodeRepo,
		showRepo,
	)

	err := episodeSearchService.ProcessSingleEpisode(int(ep.ID))
	if err != nil {
		return err
	}

	return nil
}
