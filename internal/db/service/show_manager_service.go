package service

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/tvmaze/service"
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas/dtos"
	"gorgon/pkg/services"
	"log/slog"

	"gorm.io/gorm"
)

type ShowManagerService struct {
	baseService *services.BaseService
	tvMaze      *service.TvMazeSearchService
	DB          *gorm.DB
	logger      *slog.Logger
}

func NewShowManagerService(logger *slog.Logger) *ShowManagerService {
	return &ShowManagerService{
		baseService: services.NewBaseService(),
		tvMaze:      service.NewTvMazeSearchService(logger),
		logger:      logger,
		DB:          config.GetSQLite(),
	}
}

// NOTE: Further I will need to change this (showId) to be capable to use others database beyond tvmaze
func (sm *ShowManagerService) GetEpisodes(showId int) (*[]dtos.EpisodeDto, error) {
	episodesDto, err := sm.tvMaze.SearchEpisodes(showId)
	if err != nil {
		return nil, err
	}

	return episodesDto, nil
}

func (sm *ShowManagerService) GetSeasons(showId int) (*[]dtos.SeasonDto, error) {
	seasonsDto, err := sm.tvMaze.SearchSeasons(showId)
	if err != nil {
		return nil, err
	}
	return seasonsDto, nil
}

// NOTE: Everything here will need a refactor furthermore, there're lot of problems
// We need to change the show_id on database for another better name to not confuse
// This functon call the database a lot of times
// Probably will be problems with the seasons updating.
func (sm *ShowManagerService) UpdateShowWithRelations(show *model.Show) error {
	tx := sm.DB.Begin()

	var dbShow model.Show
	if err := tx.Where("show_id = ?", show.ShowID).Find(&dbShow).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Show{}).Where("show_id = ?", show.ShowID).Updates(show).Error; err != nil {
		tx.Rollback()
		return err
	}

	var dbSeasons []model.Season
	if err := tx.Where("show_id = ?", dbShow.ID).Find(&dbSeasons).Error; err != nil {
		sm.logger.Error("Error while getting seasons from db to update seasson", slog.String("error", err.Error()))
		return nil
	}

	seasonMap := make(map[string]*model.Season)
	for i := range dbSeasons {
		e := &dbSeasons[i]
		key := fmt.Sprintf("%d:%d", e.SeasonId, e.Number)
		seasonMap[key] = e
	}

	for _, season := range show.Seasons {

		key := fmt.Sprintf("%d:%d", season.SeasonId, season.Number)
		if existing, ok := seasonMap[key]; ok {
			sm.logger.Info("Updating season", slog.Int("show_id", season.ShowId), slog.Int("season_id", season.SeasonId), slog.Int("season_number", season.Number))
			existing.SeasonId = season.SeasonId
			existing.Number = season.Number

			tx.Save(existing)
		} else {
			newSeason := model.Season{
				ShowId:   dbShow.ID,
				SeasonId: season.SeasonId,
				Number:   season.Number,
			}
			if err := tx.Create(&newSeason).Error; err != nil {
				//TODO: Error message
				return err
			}
		}
	}

	var dbEpisodes []model.Episode
	tx.Where("show_id = ?", dbShow.ID).Find(&dbEpisodes)
	episodeMap := make(map[string]*model.Episode)
	for i := range dbEpisodes {
		e := &dbEpisodes[i]
		key := fmt.Sprintf("%d:%d", e.Season, e.Number)
		episodeMap[key] = e
	}

	for _, episode := range show.Episodes {
		key := fmt.Sprintf("%d:%d", episode.Season, episode.Number)
		if existing, ok := episodeMap[key]; ok {
			sm.logger.Info("Updating episode", slog.Int("show_id", episode.ShowId), slog.Int("season", episode.Season), slog.String("episode_name", episode.Name))
			existing.Name = episode.Name
			existing.Summary = episode.Summary
			existing.Type = episode.Type
			existing.AirStamp = episode.AirStamp

			tx.Save(existing)
		} else {
			newEpisode := model.Episode{
				ShowId:   dbShow.ID,
				Name:     episode.Name,
				Summary:  episode.Summary,
				Type:     episode.Type,
				Number:   episode.Number,
				Season:   episode.Season,
				AirStamp: episode.AirStamp,
			}

			if err := tx.Create(&newEpisode).Error; err != nil {
				//TODO: Error message
				return err
			}
		}
	}

	return tx.Commit().Error
}
