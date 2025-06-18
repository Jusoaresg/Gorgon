package show

import (
	"gorgon/config"
	"gorgon/external/tvmaze/service"
	"gorgon/internal/db/model"
	"gorgon/internal/db/repository"
	"gorgon/internal/db/schema/show"
	showManager "gorgon/internal/db/service"
	"gorgon/pkg/schemas"
	"gorgon/pkg/schemas/dtos"
	"gorgon/pkg/services"
	"log/slog"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Add Show
// @Description Add Show to List, tracking must be "all", "future", or "none"
// @Tags Database/Show
// @Accept json
// @Produce json
// @Param request body show.AddShowToListRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show [post]
func AddShowToList(c echo.Context) error {
	return addShowToListHandler(c, config.GetSQLite())
}

func addShowToListHandler(c echo.Context, db *sqlx.DB) error {
	logger := config.GetLogger()
	logger.Info("Received request to Add Show To List", slog.String("endpoint", "/api/v1/database/show"), slog.String("method", "POST"))

	var request show.AddShowToListRequest

	if err := c.Bind(&request); err != nil {
		logger.Error("Failed to bind body request")
		schemas.SendError(c, 500, "Failed to bind body request")
		return err
	}

	validTrackings := map[string]bool{
		"all":    true,
		"future": true,
		"none":   true,
	}

	if !validTrackings[request.TrackingType] {
		schemas.SendError(c, 400, "Invalid tracking type: must be 'all', 'future', or 'none'")
		return echo.NewHTTPError(400, "Invalid tracking type")
	}

	tvMazeService := service.NewTvMazeSearchService(logger)
	showManagerService := showManager.NewShowManagerService(logger, db)

	idString := strconv.Itoa(request.Id)
	id64, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	showDto, err := tvMazeService.SearchByTvMazeId(id64)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	episodesDto, err := showManagerService.GetEpisodes(showDto.TvMazeID)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	services.ApplyTrackingToEpisodes(episodesDto, request.TrackingType)

	seasonsDto, err := showManagerService.GetSeasons(showDto.TvMazeID)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	show := showDto.ToModel()

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	showRepo := repository.NewShowRepository(db)
	showID, err := showRepo.CreateTx(tx, show)
	if err != nil {
		logger.Error("Failed to add show to database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to add show to database")
		return err
	}
	seasons := dtos.SeasonDtoSliceToModel(*seasonsDto, showID)
	var episodes []model.Episode

	for _, season := range seasons {
		seasonRepo := repository.NewSeasonRepository(db)
		seasonID, err := seasonRepo.CreateTx(tx, season)
		if err != nil {
			logger.Error("Failed to add season to database", slog.String("error", err.Error()))
			schemas.SendError(c, 500, "Failed to add season to database")
			return err
		}

		tempEpisodesDTO := *episodesDto
		for _, episode := range tempEpisodesDTO {
			if episode.Season == season.Number {
				ep := episode.ToModel(showID, seasonID)
				episodes = append(episodes, *ep)
			}
		}
	}

	for _, episode := range episodes {
		episodeRepo := repository.NewEpisodeRepository(db)
		if _, err := episodeRepo.CreateTx(tx, episode); err != nil {
			logger.Error("Failed to add episode to database", slog.String("error", err.Error()))
			schemas.SendError(c, 500, "Failed to add episode to database")
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to commit transaction")
		return err
	}

	schemas.SendSucess(c, "Add Show To List", &show)
	return nil
}
