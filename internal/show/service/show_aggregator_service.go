package service

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	seasonModel "github.com/jusoaresg/gorgon/internal/season/model"
	showModel "github.com/jusoaresg/gorgon/internal/show/model"
	showAliasModel "github.com/jusoaresg/gorgon/internal/show_aliases/model"

	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
)

type ShowAggregatorService struct {
	ShowRepo        showRepository.ShowRepositoryInterface
	ShowAliasesRepo showAliasRepository.ShowAliasesRepositoryInterface
	EpisodeRepo     episodeRepository.EpisodeRepositoryInterface
	SeasonRepo      seasonRepository.SeasonRepositoryInterface
}

type AggregatedShow struct {
	Show        showModel.Show
	ShowAliases []showAliasModel.ShowAlias
	Seasons     []seasonModel.Season
	Episodes    []episodeModel.Episode
}

func NewShowAggregatorService(
	showRepo showRepository.ShowRepositoryInterface,
	showAliasRepo showAliasRepository.ShowAliasesRepositoryInterface,
	episodeRepo episodeRepository.EpisodeRepositoryInterface,
	seasonRepo seasonRepository.SeasonRepositoryInterface,
) *ShowAggregatorService {
	return &ShowAggregatorService{
		ShowRepo:        showRepo,
		ShowAliasesRepo: showAliasRepo,
		EpisodeRepo:     episodeRepo,
		SeasonRepo:      seasonRepo,
	}
}

func NewShowAggregatorServiceWithDb(db *sqlx.DB) *ShowAggregatorService {
	aliasRepo := showAliasRepository.NewShowAliasesRepository(db)
	return &ShowAggregatorService{
		ShowRepo:        showRepository.NewShowRepository(db),
		ShowAliasesRepo: &aliasRepo,
		EpisodeRepo:     episodeRepository.NewEpisodeRepository(db),
		SeasonRepo:      seasonRepository.NewSeasonRepository(db),
	}
}

func (s *ShowAggregatorService) GetShowWithRelationsByTvMazeId(tvMazeID int64) (AggregatedShow, error) {

	show, err := s.ShowRepo.GetByTvMazeID(tvMazeID)
	if err != nil {
		return AggregatedShow{}, fmt.Errorf("failed to get show by tvmazeID: %w", err)
	}

	showAliases, err := s.ShowAliasesRepo.ListByShowID(show.ID)
	if err != nil {
		return AggregatedShow{}, fmt.Errorf("failed to get show aliases: %w", err)
	}

	season, err := s.SeasonRepo.ListByShowId(show.ID)
	if err != nil {
		return AggregatedShow{}, fmt.Errorf("Failed to get season by showID: %w", err)
	}

	episode, err := s.EpisodeRepo.ListByShowID(show.ID)
	if err != nil {
		return AggregatedShow{}, fmt.Errorf("Failed to get episode by showID: %w", err)
	}

	return AggregatedShow{
		Show:        show,
		ShowAliases: showAliases,
		Seasons:     season,
		Episodes:    episode,
	}, nil
}

func (s *ShowAggregatorService) GetShowWithRelationsById(id int64) (AggregatedShow, error) {

	show, err := s.ShowRepo.GetById(id)
	if err != nil {
		return AggregatedShow{}, fmt.Errorf("failed to get show by tvmazeID: %w", err)
	}

	showAliases, err := s.ShowAliasesRepo.ListByShowID(show.ID)
	if err != nil {
		return AggregatedShow{}, fmt.Errorf("failed to get show aliases: %w", err)
	}

	season, err := s.SeasonRepo.ListByShowId(show.ID)
	if err != nil {
		return AggregatedShow{}, fmt.Errorf("Failed to get season by showID: %w", err)
	}

	episode, err := s.EpisodeRepo.ListByShowID(show.ID)
	if err != nil {
		return AggregatedShow{}, fmt.Errorf("Failed to get episode by showID: %w", err)
	}

	return AggregatedShow{
		Show:        show,
		ShowAliases: showAliases,
		Seasons:     season,
		Episodes:    episode,
	}, nil
}

func (s *ShowAggregatorService) ListFullShows() ([]AggregatedShow, error) {
	shows, err := s.ShowRepo.List()
	if err != nil {
		return nil, err
	}

	var aggregated []AggregatedShow

	for _, show := range shows {
		aliases, err := s.ShowAliasesRepo.ListByShowID(show.ID)
		if err != nil {
			return nil, err
		}

		seasons, err := s.SeasonRepo.ListByShowId(show.ID)
		if err != nil {
			return nil, err
		}

		episodes, err := s.EpisodeRepo.ListByShowID(show.ID)
		if err != nil {
			return nil, err
		}

		aggregated = append(aggregated, AggregatedShow{
			Show:        show,
			ShowAliases: aliases,
			Seasons:     seasons,
			Episodes:    episodes,
		})
	}

	return aggregated, nil
}

func (s *ShowAggregatorService) ListFullShowsFiltered(search, status string) ([]AggregatedShow, error) {
	shows, err := s.ShowRepo.ListFiltered(search, status)
	if err != nil {
		return nil, err
	}

	var aggregated []AggregatedShow

	for _, show := range shows {
		aliases, err := s.ShowAliasesRepo.ListByShowID(show.ID)
		if err != nil {
			return nil, err
		}

		seasons, err := s.SeasonRepo.ListByShowId(show.ID)
		if err != nil {
			return nil, err
		}

		episodes, err := s.EpisodeRepo.ListByShowID(show.ID)
		if err != nil {
			return nil, err
		}

		aggregated = append(aggregated, AggregatedShow{
			Show:        show,
			ShowAliases: aliases,
			Seasons:     seasons,
			Episodes:    episodes,
		})
	}

	return aggregated, nil
}
