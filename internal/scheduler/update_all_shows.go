package scheduler

import (
	"gorgon/config"
	"gorgon/external/tvmaze/service"
	"gorgon/internal/db/model"
	showManager "gorgon/internal/db/service"
	"gorgon/pkg/services"
	"log/slog"
	"strconv"
	"time"
)

func UpdateAllShows() {
	logger := config.GetLogger()

	baseService := services.NewBaseService()

	var shows []model.Show
	baseService.List(&shows)

	tvMazeService := service.NewTvMazeSearchService(logger)
	showManagerService := showManager.NewShowManagerService(logger)

	apiService := services.NewAPIService("http://api.tvmaze.com", logger)

	var updates map[string]int64
	if err := apiService.Get("/updates/shows?since=day", &updates); err != nil {
		logger.Error("Failed to get updated shows from tvmaze", slog.String("error", err.Error()))
		return
	}

	for _, showOld := range shows {
		updatedAt, ok := updates[strconv.Itoa(showOld.ShowID)]
		if !ok {
			continue
		}

		if int64(showOld.Updated) < updatedAt {
			logger.Info("Updating show", slog.Int("show_id", showOld.ShowID), slog.String("title", showOld.Name))

			showModel, err := tvMazeService.SearchByTvMazeId(showOld.ShowID)
			if err != nil {
				logger.Error("Error while searching tvmaze for id while updating shows", slog.String("error", err.Error()))
				continue
			}

			episodes, err := showManagerService.GetEpisodes(showModel.ShowID)
			if err != nil {
				logger.Error("Error while getting episodes for show", slog.String("error", err.Error()))
				continue
			}

			seasons, err := showManagerService.GetSeasons(showModel.ShowID)
			if err != nil {
				logger.Error("Error while getting seasons for show", slog.String("error", err.Error()))
				continue
			}

			show := showModel.ToModel(episodes, seasons)

			if err := showManagerService.UpdateShowWithRelations(show); err != nil {
				logger.Error("Error while updating shows with relations", slog.String("error", err.Error()))
				continue
			}

			time.Sleep(5 * time.Second)
		}
	}
}
