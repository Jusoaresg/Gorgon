package front

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/jusoaresg/gorgon/internal/show/service"
	"github.com/jusoaresg/gorgon/views"
	"github.com/labstack/echo/v4"
)

type ShowsGrid struct {
	Shows   []service.AggregatedShow
	Filters struct {
		Search string
		Status string
		Sort   string
	}
}

type ShowsGridData struct {
	Data ShowsGrid
}

func (h *FrontHandler) ShowsRoute(c echo.Context) error {

	search := c.QueryParam("search")
	status := c.QueryParam("status")
	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "added"
	}

	shows, err := h.AggregatorService.ListFullShowsFiltered(search, status)
	if err != nil {
		return err
	}
	showsGrid := ShowsGrid{
		Shows: shows,
		Filters: struct {
			Search string
			Status string
			Sort   string
		}{
			Search: search,
			Status: status,
			Sort:   sort,
		},
	}

	SortShows(shows, sort)

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "shows-grid-htmx", ShowsGridData{
			Data: showsGrid,
		})
	}

	return c.Render(http.StatusOK, "layout", views.PageData{
		TemplateName: "shows",
		Data:         showsGrid,
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

func SortShows(shows []service.AggregatedShow, sortBy string) {
	switch sortBy {
	case "name":
		sort.Slice(shows, func(i, j int) bool {
			return shows[i].Show.Name < shows[j].Show.Name
		})
	case "next":
		sort.Slice(shows, func(i, j int) bool {
			return nextEpisodeTime(shows[i]) < nextEpisodeTime(shows[j])
		})
	case "added":
		sort.Slice(shows, func(i, j int) bool {
			return shows[i].Show.Updated < shows[j].Show.Updated
		})
	}
}

func nextEpisodeTime(show service.AggregatedShow) int64 {
	var next int64 = 1<<63 - 1 // maior int64 possível
	now := time.Now().Unix()

	for _, ep := range show.Episodes {
		if ep.AirStamp > now && ep.AirStamp < next {
			next = ep.AirStamp
		}
	}
	return next
}
