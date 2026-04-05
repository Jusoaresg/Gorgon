package service

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/external/prowlarr/service"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	showModel "github.com/jusoaresg/gorgon/internal/show/model"
	showService "github.com/jusoaresg/gorgon/internal/show/service"
)

type EpisodeSearcherInterface interface {
	SearchEpisodeByQuery(query string) ([]schema.SearchResponse, error)
	SearchEpisodeAliasesById(episode episodeModel.Episode, show showModel.Show, db *sqlx.DB) ([]schema.SearchResponse, error)
}

type EpisodeSearcher struct{}

func (s *EpisodeSearcher) SearchEpisodeByQuery(query string) ([]schema.SearchResponse, error) {
	logger := config.GetLogger()
	prowlarrService, err := service.NewProwlarrSearchService(logger)
	if err != nil {
		return nil, err
	}

	request := schema.SearchRequest{
		Query: query,
	}

	var response []schema.SearchResponse
	err = prowlarrService.Search(&request, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *EpisodeSearcher) SearchEpisodeAliasesById(episode episodeModel.Episode, show showModel.Show, db *sqlx.DB) ([]schema.SearchResponse, error) {
	logger := config.GetLogger()
	const cooldownDuration = 5 * time.Minute

	titleAlias, err := showService.GetNormalizedTitleAlias(show)
	if err != nil {
		logger.Error(
			"Failed to get normalized title alias",
			slog.Int64("episode_id", episode.ID),
			slog.Int64("show_id", episode.ShowID),
		)
		return nil, err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	var responses []schema.SearchResponse

	for _, title := range titleAlias {
		termKey := fmt.Sprintf("%s S%02dE%02d", title, episode.Season, episode.Number)

		if val, ok := config.ProwlarrCooldownCache.Load(termKey); ok {
			lastTime := val.(time.Time)
			if time.Since(lastTime) < cooldownDuration {
				logger.Info("Skipping search (cooldown)", slog.String("term", termKey))
				continue
			}
		}
		config.ProwlarrCooldownCache.Store(termKey, time.Now())

		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			tmpResp, err := s.SearchEpisodeByQuery(query)
			if err != nil {
				logger.Error("Search failed", slog.String("query", query), slog.String("error", err.Error()))
				return
			}
			mu.Lock()
			responses = append(responses, tmpResp...)
			mu.Unlock()
		}(termKey)
	}

	wg.Wait()
	return responses, nil
}
