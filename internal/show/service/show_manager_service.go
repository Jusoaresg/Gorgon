package service

import (
	"fmt"
	"github.com/jusoaresg/gorgon/external/tvmaze/service"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonModel "github.com/jusoaresg/gorgon/internal/season/model"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas/dtos"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type ShowManagerService struct {
	EpisodeRepo episodeRepository.EpisodeRepository
	SeasonRepo  seasonRepository.SeasonRepository
	ShowRepo    showRepository.ShowRepository
	tvMaze      *service.TvMazeSearchService
	DB          *sqlx.DB
	logger      *slog.Logger
}

func NewShowManagerService(logger *slog.Logger, db *sqlx.DB) *ShowManagerService {
	return &ShowManagerService{
		EpisodeRepo: *episodeRepository.NewEpisodeRepository(db),
		SeasonRepo:  *seasonRepository.NewSeasonRepository(db),
		ShowRepo:    *showRepository.NewShowRepository(db),
		tvMaze:      service.NewTvMazeSearchService(logger),
		logger:      logger,
		DB:          db,
	}
}

func (sm *ShowManagerService) GetEpisodes(tvMazeId int64) (*[]dtos.EpisodeDto, error) {
	episodesDto, err := sm.tvMaze.SearchEpisodes(tvMazeId)
	if err != nil {
		return nil, err
	}

	return episodesDto, nil
}

func (sm *ShowManagerService) GetSeasons(tvMazeId int64) (*[]dtos.SeasonDto, error) {
	seasonsDto, err := sm.tvMaze.SearchSeasons(tvMazeId)
	if err != nil {
		return nil, err
	}
	return seasonsDto, nil
}

// NOTE: Everything here will need a refactor furthermore, there're lot of problems
// This functon call the database a lot of times
// Probably will be problems with the seasons updating.
func (sm *ShowManagerService) UpdateShowWithRelations(showDTO dtos.ShowDto, seasonsDTO []dtos.SeasonDto, episodes []dtos.EpisodeDto) error {
	showModel, err := sm.ShowRepo.GetByTvMazeID(showDTO.TvMazeID)
	if err != nil {
		//TODO: Log message
		return err
	}

	seasonsModel, err := sm.SeasonRepo.ListByShowId(showModel.ID)
	if err != nil {
		//TODO: Better Log message
		sm.logger.Error("Error while getting seasons from db to update seasson", slog.String("error", err.Error()))
		return err
	}

	episodesModel, err := sm.EpisodeRepo.ListByShowID(showModel.ID)
	if err != nil {
		//TODO: Log message
		return err
	}

	tx, err := sm.DB.Beginx()
	if err != nil {
		return err
	}

	if err := sm.ShowRepo.UpdateTxByTvMazeID(tx, showModel); err != nil {
		tx.Rollback()
		return err
	}

	seasonMap := make(map[int]*seasonModel.Season)
	for i := range seasonsModel {
		e := &seasonsModel[i]
		seasonMap[e.Number] = e
	}

	for _, season := range seasonsModel {

		if existing, ok := seasonMap[season.Number]; ok {
			sm.logger.Info("Updating season", slog.Int64("ShowID", showModel.ID), slog.Int("SeasonNumber", season.Number))

			existing.Number = season.Number
		} else {
			newSeason := seasonModel.Season{
				ShowID: showModel.ID,
				Number: season.Number,
			}
			if _, err := sm.SeasonRepo.CreateTx(tx, newSeason); err != nil {
				//TODO: Error message
				return err
			}
		}
	}

	episodeMap := make(map[string]*episodeModel.Episode)
	for i := range episodesModel {
		e := &episodesModel[i]
		key := fmt.Sprintf("%d:%d", e.Season, e.Number)
		episodeMap[key] = e
	}

	for _, episode := range episodesModel {
		key := fmt.Sprintf("%d:%d", episode.Season, episode.Number)
		if existing, ok := episodeMap[key]; ok {
			sm.logger.Info("Updating episode", slog.Int64("show_id", episode.ShowID), slog.Int("season", episode.Season), slog.String("episode_name", episode.Name))
			existing.Name = episode.Name
			existing.Summary = episode.Summary
			existing.Type = episode.Type
			existing.AirStamp = episode.AirStamp
		} else {
			newEpisode := episodeModel.Episode{
				ShowID:   showModel.ID,
				Name:     episode.Name,
				Summary:  episode.Summary,
				Type:     episode.Type,
				Number:   episode.Number,
				Season:   episode.Season,
				AirStamp: episode.AirStamp,
			}

			_, err := sm.EpisodeRepo.CreateTx(tx, newEpisode)
			if err != nil {
				//TODO: Error message
				return err
			}
		}
	}

	return tx.Commit()
}
