package service

import (
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	seasonModel "github.com/jusoaresg/gorgon/internal/season/model"
	showModel "github.com/jusoaresg/gorgon/internal/show/model"

	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
)

type ShowAggregatorService struct {
	ShowRepo    showRepository.ShowRepositoryInterface
	EpisodeRepo episodeRepository.EpisodeRepositoryInterface
	SeasonRepo  seasonRepository.SeasonRepositoryInterface
}

type AggregatedShow struct {
	Show     showModel.Show
	Seasons  []seasonModel.Season
	Episodes []episodeModel.Episode
}

func NewShowAggregatorService(
	showRepo showRepository.ShowRepositoryInterface,
	episodeRepo episodeRepository.EpisodeRepositoryInterface,
	seasonRepo seasonRepository.SeasonRepositoryInterface,
) *ShowAggregatorService {
	return &ShowAggregatorService{
		ShowRepo:    showRepo,
		EpisodeRepo: episodeRepo,
		SeasonRepo:  seasonRepo,
	}
}

func (s *ShowAggregatorService) GetShowWithRelations(tvMazeID int64) (AggregatedShow, error) {

	show, err := s.ShowRepo.GetByTvMazeID(tvMazeID)
	if err != nil {
		return AggregatedShow{}, err
	}

	season, err := s.SeasonRepo.ListByShowId(show.ID)
	if err != nil {
		return AggregatedShow{}, err
	}

	episode, err := s.EpisodeRepo.ListByShowID(show.ID)
	if err != nil {
		return AggregatedShow{}, err
	}

	return AggregatedShow{
		Show:     show,
		Seasons:  season,
		Episodes: episode,
	}, nil
}

func (s *ShowAggregatorService) ListFullShows() ([]AggregatedShow, error) {
	shows, err := s.ShowRepo.List()
	if err != nil {
		return nil, err
	}

	var aggregated []AggregatedShow

	for _, show := range shows {
		seasons, err := s.SeasonRepo.ListByShowId(show.ID)
		if err != nil {
			return nil, err
		}

		episodes, err := s.EpisodeRepo.ListByShowID(show.ID)
		if err != nil {
			return nil, err
		}

		aggregated = append(aggregated, AggregatedShow{
			Show:     show,
			Seasons:  seasons,
			Episodes: episodes,
		})
	}

	return aggregated, nil
}
