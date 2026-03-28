package routes

import (
	"fmt"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/config"
	tvMazeService "github.com/jusoaresg/gorgon/external/tvmaze/service"
	"github.com/jusoaresg/gorgon/internal/show/service"
	"github.com/jusoaresg/gorgon/views"
	"github.com/labstack/echo/v4"
)

type FrontHandler struct {
	db                *sqlx.DB
	AggregatorService service.ShowAggregatorService
	TvMazeService     tvMazeService.TvMazeSearchService
}

func SetupFrontRouter(e *echo.Echo) {
	e.Renderer = views.NewTemplate()

	logger := config.GetLogger()

	db := config.GetSQLite()
	frontHander := &FrontHandler{
		db:                db,
		AggregatorService: *service.NewShowAggregatorServiceWithDb(db),
		TvMazeService:     *tvMazeService.NewTvMazeSearchService(logger),
	}

	staticFS, err := fs.Sub(views.FrontStaticFS, "static")
	if err != nil {
		panic(fmt.Errorf("error loading static files"))
	}
	e.StaticFS("/static", staticFS)

	e.GET("/", frontHander.ShowsRoute)

	e.GET("/show/:id", func(c echo.Context) error {
		showId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return err
		}

		db := config.GetSQLite()
		aggregatorService := service.NewShowAggregatorServiceWithDb(db)

		show, err := aggregatorService.GetShowWithRelationsById(showId)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "layout", views.PageData{
			TemplateName: "show",
			Data:         show,
			Styles:       []string{"show.css"},
		})
	})

	e.GET("/add-show", func(c echo.Context) error {
		return c.Render(http.StatusOK, "layout", views.PageData{
			TemplateName: "add-show",
			Data:         nil,
			Styles:       []string{"add-show.css"},
		})
	})

	e.GET("/add-show/:tvmaze-id/config", frontHander.AddShowConfig)

	e.GET("/calendar", func(c echo.Context) error {
		return c.Render(http.StatusOK, "layout", views.PageData{
			TemplateName: "calendar",
			Data:         nil,
			Styles:       []string{"calendar.css"},
		})
	})

	e.GET("/settings", func(c echo.Context) error {
		return c.Render(http.StatusOK, "layout", views.PageData{
			TemplateName: "config",
			Data:         nil,
			Styles:       []string{"config.css"},
		})
	})

	SetupFrontApi(e)
}
