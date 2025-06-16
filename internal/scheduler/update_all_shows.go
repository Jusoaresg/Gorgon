package scheduler

import (
	"gorgon/config"
	"gorgon/external/tvmaze/service"
	"gorgon/internal/db/repository"
	showManager "gorgon/internal/db/service"
	"gorgon/pkg/services"
	"log/slog"
	"strconv"
	"time"
)

func UpdateAllShows() {
	logger := config.GetLogger()

	showRepo := repository.NewShowRepository()
	shows, err := showRepo.List()
	if err != nil {
		return
	}

	tvMazeService := service.NewTvMazeSearchService(logger)
	showManagerService := showManager.NewShowManagerService(logger)

	apiService := services.NewAPIService("http://api.tvmaze.com", logger)

	var updates map[string]int64
	if err := apiService.Get("/updates/shows?since=day", &updates); err != nil {
		logger.Error("Failed to get updated shows from tvmaze", slog.String("error", err.Error()))
		return
	}

	for _, showOld := range shows {
		updatedAt, ok := updates[strconv.FormatInt(showOld.ID, 10)]
		if !ok {
			continue
		}

		if int64(showOld.Updated) < updatedAt {
			logger.Info("Updating show", slog.Int64("show_id", showOld.ID), slog.String("title", showOld.Name))

			showDTO, err := tvMazeService.SearchByTvMazeId(showOld.ID)
			if err != nil {
				logger.Error("Error while searching tvmaze for id while updating shows", slog.String("error", err.Error()))
				continue
			}

			episodesDTO, err := showManagerService.GetEpisodes(showDTO.TvMazeID)
			if err != nil {
				logger.Error("Error while getting episodes for show", slog.String("error", err.Error()))
				continue
			}

			seasonsDTO, err := showManagerService.GetSeasons(showDTO.TvMazeID)
			if err != nil {
				logger.Error("Error while getting seasons for show", slog.String("error", err.Error()))
				continue
			}

			if err := showManagerService.UpdateShowWithRelations(*showDTO, *seasonsDTO, *episodesDTO); err != nil {
				logger.Error("Error while updating shows with relations", slog.String("error", err.Error()))
				continue
			}

			time.Sleep(5 * time.Second)
		}
	}
}
