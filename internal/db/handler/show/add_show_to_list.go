package show

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/tvmaze/service"
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Add Show
// @Description Add Show to List
// @Tags Database/Show
// @Accept json
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show [post]
func AddShowToList(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Add Show To List", slog.String("endpoint", "/api/v1/database/show"), slog.String("method", "POST"))

	var request schemas.IdRequest

	if c.Request().Header.Get("Content-Type") != "application/json" {
		id, err := strconv.ParseInt(c.FormValue("id"), 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		request.Id = int(id)
	} else {
		if err := c.Bind(&request); err != nil {
			logger.Error("Failed to bind body request")
			schemas.SendError(c, 500, "Failed to bind body request")
			return err
		}
	}

	tvMazeService := service.NewTvMazeSearchService(logger)

	showDto, err := tvMazeService.SearchByTvMazeId(request.Id)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	episodesDto, err := tvMazeService.SearchEpisodes(request.Id)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	seasonsDto, err := tvMazeService.SearchSeasons(request.Id)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	show := model.Show{
		ShowID:    showDto.ShowID,
		Name:      showDto.Name,
		Type:      showDto.Type,
		Language:  showDto.Language,
		Status:    showDto.Status,
		Premiered: showDto.Premiered,
		Ended:     showDto.Ended,
		Rating:    showDto.Rating.Average,
		Summary:   showDto.Summary,
		Updated:   showDto.Updated,

		Seasons:  make([]model.Season, len(*seasonsDto)),
		Episodes: make([]model.Episode, len(*episodesDto)),

		Externals: model.Externals{
			Tvrage:   showDto.Externals.TvRage,
			Thetvdvb: showDto.Externals.TheTvdb,
			Imdb:     showDto.Externals.Imdb,
		},
		Image: model.Image{
			Original: showDto.Image.Original,
			Medium:   showDto.Image.Medium,
		},
	}

	for i, season := range *seasonsDto {
		show.Seasons[i] = model.Season{
			ShowId:   showDto.ShowID,
			SeasonId: season.ShowId,
			Number:   season.Number,
		}
	}

	for i, episode := range *episodesDto {
		show.Episodes[i] = model.Episode{
			ShowId:   episode.ShowId,
			Name:     episode.Name,
			Summary:  episode.Summary,
			Number:   episode.Number,
			Season:   episode.Season,
			AirStamp: episode.AirStamp,
		}
	}

	baseService := services.NewBaseService()
	if err := baseService.Add(&show); err != nil {
		logger.Error("Failed to add anime to database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to add anime to database")
		return err
	}

	schemas.SendSucess(c, "Add Show To List", &show)
	return nil
}
