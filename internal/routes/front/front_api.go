package front

import (
	"net/http"
	"strconv"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/tvmaze/service"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	seasonModel "github.com/jusoaresg/gorgon/internal/season/model"
	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	show "github.com/jusoaresg/gorgon/internal/show/handler"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showSchema "github.com/jusoaresg/gorgon/internal/show/schema"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/labstack/echo/v4"
)

func SetupFrontApi(e *echo.Echo) {
	api := e.Group("front/")

	api.POST("search-show", func(c echo.Context) error {

		var request schemas.NameRequest
		if err := c.Bind(&request); err != nil {
			return nil
		}

		logger := config.GetLogger()
		db := config.GetSQLite()

		tvMazeService := service.NewTvMazeSearchService(logger)
		showRepo := showRepository.NewShowRepository(db)

		showManager := service.NewShowManager(*tvMazeService, *showRepo, logger)

		shows, err := showManager.SearchAndEnrich(request.Name)
		if err != nil {
			return err
		}

		data := map[string]any{
			"Shows": shows,
		}

		return c.Render(http.StatusOK, "add-show-card", data)
	})

	api.POST("add-show", AddShow)

	api.GET("show/:id/modal/edit", EditShowModal)

	api.GET("episode/:id/modal/tracking", ChangeEpisodeTracking)
	api.GET("season/:id/modal/tracking", ChangeSeasonTracking)
}

func AddShow(c echo.Context) error {
	logger := config.GetLogger()

	id, err := strconv.ParseInt(c.FormValue("id"), 10, 64)
	if err != nil {
		return err
	}

	trackingType := c.FormValue("monitor")

	request := showSchema.AddShowToListRequest{
		Id:           int(id),
		TrackingType: trackingType,
	}
	show.AddShowToListHandler(c, &request, config.GetSQLite(), logger)

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func EditShowModal(c echo.Context) error {
	id := c.Param("id")
	return c.Render(http.StatusOK, "edit-show-modal", id)
}

func ChangeEpisodeTracking(c echo.Context) error {
	epIdStr := c.Param("id")
	epIdInt, err := strconv.Atoi(epIdStr)
	if err != nil {
		return err
	}

	epRepo := episodeRepository.NewEpisodeRepository(config.GetSQLite())
	episode, err := epRepo.GetByID(int64(epIdInt))
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "episode-tracking-modal", episode)
}

func ChangeSeasonTracking(c echo.Context) error {
	seasonIdStr := c.Param("id")
	seasonIdInt, err := strconv.Atoi(seasonIdStr)
	if err != nil {
		return err
	}

	db := config.GetSQLite()

	seasonRepo := seasonRepository.NewSeasonRepository(db)
	episodeRepo := episodeRepository.NewEpisodeRepository(db)

	episodes, err := episodeRepo.ListBySeasonID(seasonIdInt)
	if err != nil {
		return err
	}

	season, err := seasonRepo.GetById(int64(seasonIdInt))
	if err != nil {
		return err
	}

	type seasonModal struct {
		Season     seasonModel.Season
		EpisodeIds []int
	}

	episodesIds := make([]int, 0, len(episodes))
	for _, e := range episodes {
		episodesIds = append(episodesIds, int(e.ID))
	}
	return c.Render(http.StatusOK, "season-tracking-modal", seasonModal{
		Season:     season,
		EpisodeIds: episodesIds,
	})
}
