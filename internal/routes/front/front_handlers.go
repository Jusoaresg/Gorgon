package front

import (
	"net/http"
	"strconv"

	"github.com/jusoaresg/gorgon/views"
	"github.com/labstack/echo/v4"
)

func (h *FrontHandler) ShowsRoute(c echo.Context) error {

	shows, err := h.AggregatorService.ListFullShows()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "layout", views.PageData{
		TemplateName: "shows",
		Data:         shows,
		Styles:       []string{"shows.css"},
	})
}

func (h *FrontHandler) AddShowConfig(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("tvmaze-id"), 10, 64)
	if err != nil {
		return err
	}

	show, err := h.TvMazeService.SearchByTvMazeId(id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "layout", views.PageData{
		TemplateName: "add-show-config",
		Data:         show,
		Styles:       []string{"add-show-config.css"},
	})
}
