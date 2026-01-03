package show

import (
	"log/slog"
	"strconv"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/tvmaze/service"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	"github.com/jusoaresg/gorgon/internal/show/model"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/internal/show/schema"
	showManager "github.com/jusoaresg/gorgon/internal/show/service"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/jusoaresg/gorgon/pkg/schemas/dtos"
	"github.com/jusoaresg/gorgon/pkg/services"

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
	logger := config.GetLogger()
	logger.Info("Received request to Add Show To List", slog.String("endpoint", "/api/v1/database/show"), slog.String("method", "POST"))

	var request show.AddShowToListRequest

	if err := c.Bind(&request); err != nil {
		logger.Error("Failed to bind body request")
		schemas.SendError(c, 500, "Failed to bind body request")
		return err
	}

	show, err := AddShowToListHandler(c, &request, config.GetSQLite(), logger)
	if err != nil {
		return err
	}

	schemas.SendSuccess(c, "Add Show To List", &show)
	return nil
}

func AddShowToListHandler(c echo.Context, request *show.AddShowToListRequest, db *sqlx.DB, logger *slog.Logger) (*model.Show, error) {

	validTrackings := map[string]bool{
		"all":    true,
		"future": true,
		"none":   true,
	}

	if !validTrackings[request.TrackingType] {
		schemas.SendError(c, 400, "Invalid tracking type: must be 'all', 'future', or 'none'")
		return nil, echo.NewHTTPError(400, "Invalid tracking type")
	}

	tvMazeService := service.NewTvMazeSearchService(logger)
	showManagerService := showManager.NewShowManagerService(logger, db)

	idString := strconv.Itoa(request.Id)
	id64, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return nil, err
	}

	showDto, err := tvMazeService.SearchByTvMazeId(id64)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return nil, err
	}

	episodesDto, err := showManagerService.GetEpisodes(showDto.TvMazeID)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return nil, err
	}

	services.ApplyTrackingToEpisodes(episodesDto, request.TrackingType)

	seasonsDto, err := showManagerService.GetSeasons(showDto.TvMazeID)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return nil, err
	}

	show := showDto.ToModel()

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}

	showRepo := showRepository.NewShowRepository(db)
	showID, err := showRepo.CreateTx(tx, show)
	if err != nil {
		logger.Error("Failed to add show to database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to add show to database")
		return nil, err
	}

	aliasRepo := showAliasRepository.NewShowAliasesRepository(db)
	aliases := showDto.ToAliasModel()
	for _, alias := range aliases {
		alias.ShowID = showID
		alias.Source = "tvmaze"
		_, err := aliasRepo.CreateTx(tx, alias)
		if err != nil {
			logger.Error(
				"failed to create show alias",
				slog.String("error", err.Error()),
				slog.String("alias", alias.Alias),
				slog.String("country", alias.Country),
			)
			continue
		}
	}

	seasons := dtos.SeasonDtoSliceToModel(*seasonsDto, showID)
	var episodes []episodeModel.Episode

	for _, season := range seasons {
		seasonRepo := seasonRepository.NewSeasonRepository(db)
		seasonID, err := seasonRepo.CreateTx(tx, season)
		if err != nil {
			logger.Error("Failed to add season to database", slog.String("error", err.Error()))
			schemas.SendError(c, 500, "Failed to add season to database")
			return nil, err
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
		episodeRepo := episodeRepository.NewEpisodeRepository(db)
		if _, err := episodeRepo.CreateTx(tx, episode); err != nil {
			logger.Error("Failed to add episode to database", slog.String("error", err.Error()))
			schemas.SendError(c, 500, "Failed to add episode to database")
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("Failed to commit transaction", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to commit transaction")
		return nil, err
	}

	return &show, nil
}
