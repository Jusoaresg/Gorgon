package service

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/external/tvmaze/schema"
	"github.com/jusoaresg/gorgon/internal/show/repository"
)

type ShowManager struct {
	tvMaze   TvMazeSearchService
	showRepo repository.ShowRepository
	logger   *slog.Logger
}

func NewShowManager(tvMaze TvMazeSearchService, showRepo repository.ShowRepository, logger *slog.Logger) *ShowManager {
	return &ShowManager{
		tvMaze:   tvMaze,
		showRepo: showRepo,
		logger:   logger,
	}
}

func (s *ShowManager) SearchAndEnrich(name string) (*[]schema.SearchResult, error) {
	response, err := s.tvMaze.SearchByName(name)
	if err != nil {
		s.logger.Error("Error while searching for name", slog.String("error", err.Error()))
		return nil, err
	}

	existingShows, err := s.showRepo.List()
	if err != nil {
		return nil, err
	}

	existingsMap := make(map[int64]bool)
	for _, s := range existingShows {
		existingsMap[s.TvMazeID] = true
	}

	var enriched []schema.SearchResult
	for _, r := range *response {
		enriched = append(enriched, schema.SearchResult{
			Show:    r.Show,
			IsAdded: existingsMap[r.Show.TvMazeID],
		})
	}
	return &enriched, nil
}
