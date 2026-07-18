package service

import (
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/external/prowlarr/service"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	showModel "github.com/jusoaresg/gorgon/internal/show/model"
	showService "github.com/jusoaresg/gorgon/internal/show/service"
)

var latinRegex = regexp.MustCompile(`^[a-zA-Z0-9\s\-_!?',.:]+$`)

type EpisodeSearcherInterface interface {
	SearchEpisodeByQuery(query string) ([]schema.SearchResponse, error)
	SearchEpisodeAliasesById(episode episodeModel.Episode, show showModel.Show, db *sqlx.DB) ([]schema.SearchResponse, error)
}

type EpisodeSearcher struct{}

func (s *EpisodeSearcher) SearchEpisodeByQuery(query string) ([]schema.SearchResponse, error) {
	logger := config.GetLogger()
	prowlarrIndexerService := service.NewProwlarrIndexerService(logger)

	var indexers []schema.IndexerResponse
	err := prowlarrIndexerService.GetIndexers(&indexers)
	if err != nil {
		return nil, err
	}
	var indexerIds []int
	for _, indexer := range indexers {
		if indexer.Enabled {
			indexerIds = append(indexerIds, indexer.Id)
		}
	}

	prowlarrService, err := service.NewProwlarrSearchService(logger)
	if err != nil {
		return nil, err
	}

	request := schema.SearchRequest{
		Query: query,
	}

	var response []schema.SearchResponse
	err = prowlarrService.Search(&request, &response, indexerIds...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *EpisodeSearcher) SearchEpisodeAliasesById(episode episodeModel.Episode, show showModel.Show, db *sqlx.DB) ([]schema.SearchResponse, error) {
	logger := config.GetLogger()
	const cooldownDuration = 5 * time.Minute

	rawAliases, err := showService.GetNormalizedTitleAlias(show)
	if err != nil {
		logger.Error(
			"Failed to get normalized title alias",
			slog.Int64("episode_id", episode.ID),
			slog.Int64("show_id", episode.ShowID),
		)
		return nil, err
	}

	var cleanAliases []string
	for _, alias := range rawAliases {
		if alias == show.Name {
			continue
		}

		if latinRegex.MatchString(alias) {
			cleanAliases = append(cleanAliases, alias)
		}
	}

	var allResponses []schema.SearchResponse

	for _, title := range cleanAliases {
		termKey := fmt.Sprintf("%s S%02dE%02d", title, episode.Season, episode.Number)

		if val, ok := config.ProwlarrCooldownCache.Load(termKey); ok {
			lastTime := val.(time.Time)
			if time.Since(lastTime) < cooldownDuration {
				logger.Info("Skipping search (cooldown)", slog.String("term", termKey))
				continue
			}
		}
		config.ProwlarrCooldownCache.Store(termKey, time.Now())

		tmpResp, err := s.SearchEpisodeByQuery(termKey)
		if err != nil {
			logger.Error(
				"Search failed",
				slog.String("query", termKey),
				slog.String("error", err.Error()),
			)
			continue
		}

		if len(tmpResp) <= 0 {
			continue
		}

		allResponses = append(allResponses, tmpResp...)

		if s.hasGoodQuality(tmpResp) {
			logger.Info("Good quality torrent found! Stopping extra searches", slog.String("used_alias", title))
			break
		}
	}

	return allResponses, nil
}

func (s *EpisodeSearcher) hasGoodQuality(results []schema.SearchResponse) bool {
	for _, response := range results {
		score := ScoreEpisode(response)
		if score >= 60 {
			return true
		}
	}
	return false
}
