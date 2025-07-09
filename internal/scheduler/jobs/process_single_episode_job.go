package jobs

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	prowlarr "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrent "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	EpisodeService "github.com/jusoaresg/gorgon/internal/episode/service"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showService "github.com/jusoaresg/gorgon/internal/show/service"
)

var cooldownCache sync.Map

const cooldownDuration = 5 * time.Minute

func init() {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		for range ticker.C {
			now := time.Now()
			cooldownCache.Range(func(key, value any) bool {
				if t, ok := value.(time.Time); ok && now.Sub(t) > cooldownDuration*2 {
					cooldownCache.Delete(key)
				}
				return true
			})
		}
	}()
}

func ProcessBackfillSingleEpisode(ep *model.Episode, prowlarrService *prowlarr.ProwlarrSearchService, qbittorrentService *qbittorrent.QBittorrentService) error {
	logger := config.GetLogger().WithGroup("jobs").With("name", "ProcessSingleEpisode")
	db := config.GetSQLite()

	showRepo := showRepository.NewShowRepository(db)

	show, err := showRepo.GetById(ep.ShowID)
	if err != nil {
		return err
	}

	if aired := ep.HasAired(); !aired {
		logger.Warn(
			"episode has not aired yet",
			slog.Int64("episode_id", ep.ID),
			slog.Int64("show_id", ep.ShowID),
		)
		return nil
	}
	logger.Info("searching if episode is avaible", slog.String("show_name", show.Name), slog.Int64("show_id", show.ID), slog.Int("episode", ep.Number), slog.String("episode_name", ep.Name), slog.String("tracking", string(ep.Tracking)))

	titleAlias, err := showService.GetNormalizedTitleAlias(show)
	if err != nil {
		logger.Error(
			"failed to get normalized title alias",
			slog.Int64("episode_id", ep.ID),
			slog.Int64("show_id", ep.ShowID),
		)
		return err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	var response []schema.SearchResponse
	for _, title := range titleAlias {
		termKey := fmt.Sprintf("%s S%02dE%02d", title, ep.Season, ep.Number)

		mu.Lock()
		if val, ok := cooldownCache.Load(termKey); ok {
			lastTime := val.(time.Time)
			if time.Since(lastTime) < cooldownDuration {
				logger.Info("skipping search (cooldown)", slog.String("term", termKey))
				mu.Unlock()
				continue
			}
		}

		cooldownCache.Store(termKey, time.Now())
		mu.Unlock()

		wg.Add(1)
		go func(t, key string) {
			defer wg.Done()
			tmpResponse, err := EpisodeService.SearchEpisode(key)
			if err != nil {
				return
			}
			mu.Lock()
			response = append(response, tmpResponse...)
			mu.Unlock()
		}(title, termKey)

	}

	wg.Wait()

	if len(response) == 0 {
		logger.Info("no avaible episode found", slog.String("show", show.Name), slog.Int("episode", ep.Number))
		return nil
	}

	response = EpisodeService.FilterRequiredWords(response)
	if len(response) <= 0 {
		logger.Info("no avaible episode found", slog.String("show", show.Name), slog.Int("episode", ep.Number))
	}

	response = EpisodeService.FilterByEpisodeScore(response)
	if len(response) <= 0 {
		logger.Info(
			"no available episodes found after filtering by episode score",
			slog.Int64("episode_id", ep.ID),
			slog.Int64("show_id", ep.ShowID),
		)
		return nil
	}

	logger.Debug("final chosen torrent", slog.String("filename", response[0].Filename))
	if err := EpisodeService.DownloadEpisode(*ep, response[0]); err != nil {
		return err
	}

	logger.Info("added torrent to qbittorrent", slog.String("show", show.Name), slog.Int("episode", ep.Number))

	return nil
}
