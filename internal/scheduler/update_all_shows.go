package scheduler

import (
	"gorgon/config"
	"gorgon/external/tvmaze/service"
	"gorgon/internal/db/model"
	showManager "gorgon/internal/db/service"
	"gorgon/pkg/services"
	"log/slog"
	"time"
)

func UpdateAllShows() {
	logger := config.GetLogger()

	baseService := services.NewBaseService()

	var shows []model.Show
	baseService.List(&shows)

	tvMazeService := service.NewTvMazeSearchService(logger)
	showManagerService := showManager.NewShowManagerService(logger)
	for _, showOld := range shows {

		showModel, err := tvMazeService.SearchByTvMazeId(showOld.ShowID)
		if err != nil {
			logger.Error("Error while searching tvmaze for id while updating shows", slog.String("error", err.Error()))
			return
		}

		episodes, err := showManagerService.GetEpisodes(showModel.ShowID)
		seasons, err := showManagerService.GetSeasons(showModel.ShowID)

		show := showModel.ToModel(episodes, seasons)

		if err := showManagerService.UpdateShowWithRelations(show); err != nil {
			logger.Error("Error while updating shows with relations", slog.String("error", err.Error()))
			return
		}

		time.Sleep(30 * time.Second)
	}

}
