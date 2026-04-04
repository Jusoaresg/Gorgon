package show

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/config"
	tvMazeService "github.com/jusoaresg/gorgon/external/tvmaze/service"
	"github.com/jusoaresg/gorgon/internal/show/repository"
	showService "github.com/jusoaresg/gorgon/internal/show/service"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/jusoaresg/gorgon/pkg/schemas/dtos"

	"github.com/labstack/echo/v4"
)

type UpdateShowData struct {
	UpdateShow   dtos.ShowDto `json:"updatedShow"`
	ToastMessage string       `json:"toastMessage"`
}

// @BasePath /api/v1

// @Summary Update Show Info
// @Description Updates the Show info from TvMze
// @Tags Database/Show
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/update-info [post]
func UpdateShowInfo(c echo.Context) error {
	return updateShowInfoHandler(c, config.GetSQLite())
}

func updateShowInfoHandler(c echo.Context, db *sqlx.DB) error {
	logger := config.GetLogger()
	logger.Info("Received request to update show info", slog.String("endpoint", "/database/show/update-info"), slog.String("method", "POST"))

	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		return err
	}
	showRepo := repository.NewShowRepository(db)
	tvMazeService := tvMazeService.NewTvMazeSearchService(logger)
	showManager := showService.NewShowManagerService(logger, db)

	updatedShow, err := UpdateSingleShowInfo(showRepo, tvMazeService, showManager, logger, request.Id)
	if err != nil {
		logger.Error("Failed to update show information", slog.Int("show_id", int(request.Id)))
		schemas.SendError(c, 500, "Failed to update show info", UpdateShowData{
			UpdateShow:   *updatedShow,
			ToastMessage: "Failed to update show info",
		})
	}

	logger.Info("Show info updated successfully")
	schemas.SendSuccess(c, "Update Show Info", UpdateShowData{
		UpdateShow:   *updatedShow,
		ToastMessage: "Show info updated",
	})
	return nil
}

func UpdateSingleShowInfo(
	showRepo *repository.ShowRepository,
	tvMazeService *tvMazeService.TvMazeSearchService,
	showManager *showService.ShowManagerService,
	logger *slog.Logger,
	showId int64,
) (*dtos.ShowDto, error) {

	show, err := showRepo.GetById(showId)
	if err != nil {
		logger.Error("Failed to fetch show from repository", slog.Int64("show_id", showId), slog.String("error", err.Error()))
		return nil, err
	}

	showDTO, err := tvMazeService.SearchByTvMazeId(show.TvMazeID)
	if err != nil {
		logger.Error("Failed to fetch show info from TVMaze", slog.Int64("tvmaze_id", show.TvMazeID), slog.String("error", err.Error()))
		return nil, err
	}

	episodesDTO, err := showManager.GetEpisodes(show.TvMazeID)
	if err != nil {
		logger.Error("Failed to fetch episodes from TVMaze", slog.Int64("tvmaze_id", show.TvMazeID), slog.String("error", err.Error()))
		return nil, err
	}

	seasonsDTO, err := showManager.GetSeasons(show.TvMazeID)
	if err != nil {
		logger.Error("Failed to fetch seasons from TVMaze", slog.Int64("tvmaze_id", show.TvMazeID), slog.String("error", err.Error()))
		return nil, err
	}

	err = showManager.UpdateShowWithRelations(*showDTO, *seasonsDTO, *episodesDTO)
	if err != nil {
		logger.Error("Failed to update show info", slog.Int64("tvmaze_id", show.TvMazeID), slog.String("error", err.Error()))
		return nil, err
	}

	return showDTO, nil
}
