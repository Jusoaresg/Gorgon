package scheduler

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/tvmaze/service"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showManager "github.com/jusoaresg/gorgon/internal/show/service"
	"github.com/jusoaresg/gorgon/pkg/services"
)

func UpdateAllShows() {
	logger := config.GetLogger().WithGroup("scheduler").With("name", "UpdateAllShows")
	db := config.GetSQLite()

	showRepo := showRepository.NewShowRepository(db)
	shows, err := showRepo.List()
	if err != nil {
		logger.Error("failed to list shows from repository", slog.String("error", err.Error()))
		return
	}

	tvMazeService := service.NewTvMazeSearchService(logger)
	showManagerService := showManager.NewShowManagerService(logger, db)
	apiService := services.NewAPIService("http://api.tvmaze.com", logger)

	var updates map[string]int64

	url := "/updates/shows?since=week"
	if err := apiService.Get(url, &updates); err != nil {
		logger.Error("failed to get updated shows from tvmaze", slog.String("error", err.Error()))
		return
	}

	updatedCount := 0
	for _, showOld := range shows {
		updatedAt, ok := updates[strconv.FormatInt(showOld.TvMazeID, 10)]
		if !ok {
			continue
		}

		if int64(showOld.Updated) < updatedAt {
			logger.Info(
				"updating show",
				slog.Int64("show_id", showOld.ID),
				slog.Int64("tv_maze_id", showOld.TvMazeID),
				slog.String("title", showOld.Name),
			)

			showDTO, err := tvMazeService.SearchByTvMazeId(showOld.TvMazeID)
			if err != nil {
				logger.Error("error while searching tvmaze for id while updating shows", slog.String("error", err.Error()))
				continue
			}

			episodesDTO, err := showManagerService.GetEpisodes(showDTO.TvMazeID)
			if err != nil {
				logger.Error("error while getting episodes for show", slog.String("error", err.Error()))
				continue
			}

			seasonsDTO, err := showManagerService.GetSeasons(showDTO.TvMazeID)
			if err != nil {
				logger.Error("error while getting seasons for show", slog.String("error", err.Error()))
				continue
			}

			if err := showManagerService.UpdateShowWithRelations(*showDTO, *seasonsDTO, *episodesDTO); err != nil {
				logger.Error("error while updating shows with relations", slog.String("error", err.Error()))
				continue
			}

			updatedCount++
			time.Sleep(15 * time.Second)
		}

		logger.Info("finished updating shows", slog.Int("total_updated", updatedCount))
	}
}
