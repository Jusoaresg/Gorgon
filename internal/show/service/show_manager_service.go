package service

import (
	"fmt"
	"log/slog"

	"github.com/jusoaresg/gorgon/external/tvmaze/service"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonModel "github.com/jusoaresg/gorgon/internal/season/model"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas/dtos"

	"github.com/jmoiron/sqlx"
)

type ShowManagerService struct {
	ShowAggregator ShowAggregatorService
	ShowRepo       showRepository.ShowRepositoryInterface
	SeasonRepo     seasonRepository.SeasonRepositoryInterface
	EpisodeRepo    episodeRepository.EpisodeRepositoryInterface
	tvMaze         *service.TvMazeSearchService
	DB             *sqlx.DB
	logger         *slog.Logger
}

func NewShowManagerService(logger *slog.Logger, db *sqlx.DB) *ShowManagerService {
	showRepo := showRepository.NewShowRepository(db)
	seasonRepo := seasonRepository.NewSeasonRepository(db)
	episodeRepo := episodeRepository.NewEpisodeRepository(db)

	return &ShowManagerService{
		ShowAggregator: *NewShowAggregatorService(
			showRepo,
			episodeRepo,
			seasonRepo,
		),
		ShowRepo:    showRepo,
		SeasonRepo:  seasonRepo,
		EpisodeRepo: episodeRepo,
		tvMaze:      service.NewTvMazeSearchService(logger),
		logger:      logger,
		DB:          db,
	}
}

// TODO: Move this two Get to another file
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

// NOTE: Everything here will need a refactor later, there're lot of problems
// Probably will be problems with the seasons updating.
func (sm *ShowManagerService) UpdateShowWithRelations(showDTO dtos.ShowDto, seasonsDTO []dtos.SeasonDto, episodes []dtos.EpisodeDto) error {

	aggregatedShow, err := sm.ShowAggregator.GetShowWithRelations(showDTO.TvMazeID)
	if err != nil {
		sm.logger.Error("Error while getting seasons from db to update seasson", slog.String("error", err.Error()))
		return err
	}

	showModel := aggregatedShow.Show
	episodesModel := aggregatedShow.Episodes
	seasonsModel := aggregatedShow.Seasons

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

	for _, season := range seasonsDTO {

		if existing, ok := seasonMap[season.Number]; ok {
			sm.logger.Info("Updating season", slog.Int64("ShowID", showModel.ID), slog.Int("SeasonNumber", season.Number))

			existing.Number = season.Number
		} else {
			newSeason := seasonModel.Season{
				ShowID: showModel.ID,
				Number: season.Number,
			}
			createdSeasonID, err := sm.SeasonRepo.CreateTx(tx, newSeason)
			if err != nil {
				//TODO: Error message
				return err
			}
			seasonMap[season.Number] = &seasonModel.Season{
				ID:     createdSeasonID,
				ShowID: showModel.ID,
				Number: season.Number,
			}
		}
	}

	episodeMap := make(map[string]*episodeModel.Episode)
	for i := range episodesModel {
		e := &episodesModel[i]
		key := fmt.Sprintf("%d:%d", e.Season, e.Number)
		episodeMap[key] = e
	}

	for _, episode := range episodes {
		key := fmt.Sprintf("%d:%d", episode.Season, episode.Number)
		if existing, ok := episodeMap[key]; ok {
			sm.logger.Info("Updating episode", slog.Int64("show_id", showModel.ID), slog.Int("season", episode.Season), slog.String("episode_name", episode.Name))
			existing.Name = episode.Name
			existing.Summary = episode.Summary
			existing.Type = episode.Type
			existing.AirStamp = episode.AirStamp

			if err := sm.EpisodeRepo.UpdateTx(tx, *existing); err != nil {
				return err
			}
		} else {
			season := seasonMap[episode.Season]
			if season == nil {
				sm.logger.Error("Missing season for episode", slog.Int("season_number", episode.Season))
				return fmt.Errorf("season %d not found when creating episode", episode.Season)
			}

			newEpisode := episodeModel.Episode{
				ShowID:   showModel.ID,
				SeasonID: season.ID,
				Name:     episode.Name,
				Summary:  episode.Summary,
				Type:     episode.Type,
				Number:   episode.Number,
				Season:   episode.Season,
				AirStamp: episode.AirStamp,
				Tracking: model.TrackingWanted,
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
