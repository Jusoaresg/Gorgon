package web

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jusoaresg/gorgon/config"
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

func (h *Handler) ShowsRoute(c echo.Context) error {
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

	return views.Render(c, views.View{
		Layout:    "layout",
		Default:   "shows",
		Templates: map[string]string{"grid": "shows-grid-htmx"},
		Data:      showsGrid,
		Styles:    []string{"shows.css"},
	})
}

func (h *Handler) ShowRoute(c echo.Context) error {
	showId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	show, err := h.AggregatorService.GetShowWithRelationsById(showId)
	if err != nil {
		return err
	}

	return views.Render(c, views.View{
		Layout:  "layout",
		Default: "show",
		Data:    show,
		Styles:  []string{"show.css"},
	})
}

func (h *Handler) AddShowRoute(c echo.Context) error {
	return views.Render(c, views.View{
		Layout:  "layout",
		Default: "add-show",
		Data:    nil,
		Styles:  []string{"add-show.css"},
	})
}

func (h *Handler) AddShowConfigRoute(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("tvmaze-id"), 10, 64)
	if err != nil {
		return err
	}

	show, err := h.TvMazeService.SearchByTvMazeId(id)
	if err != nil {
		return err
	}

	return views.Render(c, views.View{
		Layout:  "layout",
		Default: "add-show-config",
		Data:    show,
		Styles:  []string{"add-show-config.css"},
	})
}

func (h *Handler) CalendarRoute(c echo.Context) error {
	return views.Render(c, views.View{
		Layout:  "layout",
		Default: "calendar",
		Data:    nil,
		Styles:  []string{"calendar.css"},
	})
}

type SettingType int

const (
	GorgonSettings SettingType = iota
	ProwlarrSettings
	TorrentSettings
)

var SettingTypes = map[string]SettingType{
	"gorgon":   GorgonSettings,
	"prowlarr": ProwlarrSettings,
	"torrent":  TorrentSettings,
}

type SettingsData struct {
	Type       SettingType
	TypeString string
	Settings   any
}

func (h *Handler) SettingsRoute(c echo.Context) error {
	return views.Render(c, views.View{
		Layout:  "layout",
		Default: "settings",
		Data: SettingsData{
			Type:       GorgonSettings,
			TypeString: strings.Join([]string{"gorgon", "Settings"}, ""),
		},
		Styles: []string{"settings.css"},
	})
}

func (h *Handler) SettingsTypeRoute(c echo.Context) error {
	settingsType := c.Param("type")

	settingType, ok := SettingTypes[settingsType]
	if !ok {
		settingType = GorgonSettings
		settingsType = "gorgon"
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return views.Render(c, views.View{
			Layout:  "layout",
			Default: "settings",
			Data:    nil,
			Styles:  []string{"settings.css"},
		})
	}

	return views.Render(c, views.View{
		Layout:  "layout",
		Default: "settings",
		Data: SettingsData{
			Type:       settingType,
			TypeString: strings.Join([]string{settingsType, "Settings"}, ""),
			Settings:   cfg,
		},
		Styles: []string{"settings.css"},
	})
}

// INFO: Helpers
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
