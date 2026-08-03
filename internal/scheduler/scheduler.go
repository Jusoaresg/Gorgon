package scheduler

import (
	"errors"
	"log/slog"
	"time"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/scheduler/workers"

	prowlarrService "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrentService "github.com/jusoaresg/gorgon/external/qbittorrent/service"
)

func Start() {
	logger := config.GetLogger()

	go func() {
		var prowlarr *prowlarrService.ProwlarrSearchService
		var qbittorrent *qbittorrentService.QBittorrentService
		var err error

		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()

		for ; true; <-ticker.C {

			if prowlarr == nil {
				prowlarr, err = prowlarrService.NewProwlarrSearchService(logger)
				if err == nil {
					go workers.StartRssEpisodeFetcherWorker(5, prowlarr)
					logger.Info("Prowlarr connected and RSS Worker initiated")
				} else {
					var cfgErr prowlarrService.ErrProwlarrHostPortNotSet
					if errors.As(err, &cfgErr) && cfgErr.IsProwlarrConfigWarn() {
						logger.Warn(err.Error())
					} else {
						logger.Error("Failed to create prowlarr service", slog.String("error", err.Error()))
					}
				}
			}

			if qbittorrent == nil {
				qbittorrent, err = qbittorrentService.NewQBittorrentService(logger)
				if err == nil {
					go workers.VerifySnatchedDownloadsWorker(5, qbittorrent)
					logger.Info("QBittorrent connected and Verify Worker initiated")
				} else {
					var cfgErr qbittorrentService.ErrQBittorrentHostPortNotSet
					if errors.As(err, &cfgErr) && cfgErr.IsProwlarrConfigWarn() {
						logger.Warn(err.Error())
					} else {
						logger.Error("Failed to create qbittorrent service", slog.String("error", err.Error()))
					}
				}
			}

			if prowlarr != nil && qbittorrent != nil {
				go workers.StartEpisodeSearchWorker(2, prowlarr, qbittorrent)
				logger.Info("All services connected. Scheduler completely operational")
				return
			}
		}
	}()
}
