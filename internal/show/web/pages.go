package web

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func computeWeekStart(weekParam string) time.Time {
	now := time.Now().UTC()
	if weekParam != "" {
		parsed, err := time.Parse("2006-01-02", weekParam)
		if err == nil {
			return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
		}
	}
	return mondayOfWeek(now)
}

func mondayOfWeek(now time.Time) time.Time {
	daysSinceMonday := int(now.Weekday()+6) % 7
	monday := now.AddDate(0, 0, -daysSinceMonday)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
}

func (h *Handler) computeCalendarData(weekParam string) (CalendarData, error) {
	weekStart := computeWeekStart(weekParam)
	weekEnd := weekStart.AddDate(0, 0, 7)
	today := time.Now().UTC().Format("2006-01-02")

	var episodes []CalendarEpisode
	err := h.DB.Select(&episodes, `
		SELECT e.id, e.show_id, e.name, e.number, e.season, e.airstamp, e.tracking,
		       s.name AS show_name, s.image_medium AS show_image, s.status AS show_status
		FROM episodes e
		JOIN shows s ON e.show_id = s.id
		WHERE e.airstamp >= ? AND e.airstamp < ?
		ORDER BY e.airstamp ASC
	`, weekStart.Unix(), weekEnd.Unix())
	if err != nil {
		return CalendarData{}, err
	}

	days := make([]CalendarDay, 7)
	for i := range 7 {
		day := weekStart.AddDate(0, 0, i)
		days[i] = CalendarDay{
			Date:        day.Format("2006-01-02"),
			DayName:     day.Format("Monday"),
			DisplayDate: day.Format("2 Jan"),
			Episodes:    []CalendarEpisode{},
		}
	}

	for _, ep := range episodes {
		t := time.Unix(ep.AirStamp, 0).UTC()
		dayIdx := int(t.Sub(weekStart) / (24 * time.Hour))
		if dayIdx >= 0 && dayIdx < 7 {
			days[dayIdx].Episodes = append(days[dayIdx].Episodes, ep)
		}
	}

	currentWeekStart := computeWeekStart("")
	isCurrentWeek := weekStart.Equal(currentWeekStart)

	return CalendarData{
		Days:          days,
		WeekStart:     weekStart.Format("Jan 2"),
		WeekEnd:       weekEnd.AddDate(0, 0, -1).Format("Jan 2, 2006"),
		PrevWeek:      weekStart.AddDate(0, 0, -7).Format("2006-01-02"),
		NextWeek:      weekStart.AddDate(0, 0, 7).Format("2006-01-02"),
		Today:         today,
		IsCurrentWeek: isCurrentWeek,
	}, nil
}

func (h *Handler) CalendarRoute(c echo.Context) error {
	data, err := h.computeCalendarData(c.QueryParam("week"))
	if err != nil {
		return err
	}

	return views.Render(c, views.View{
		Layout:    "layout",
		Default:   "calendar",
		Templates: map[string]string{"week": "calendar-week"},
		Data:      data,
		Styles:    []string{"calendar.css"},
	})
}

func (h *Handler) CalendarWeekHTMX(c echo.Context) error {
	data, err := h.computeCalendarData(c.QueryParam("week"))
	if err != nil {
		return err
	}

	return c.Render(200, "calendar-week", views.PageData{Data: data})
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

type LogEntry struct {
	Time    string
	Level   string
	Message string
	Context string
}

type LogFileInfo struct {
	Name      string
	Size      string
	IsCurrent bool
}

type LogsData struct {
	Files       []LogFileInfo
	Entries     []LogEntry
	CurrentFile string
	Search      string
	Level       string
}

func (h *Handler) LogsRoute(c echo.Context) error {
	logsPath := config.LogsPath

	search := c.QueryParam("search")
	levelFilter := c.QueryParam("level")
	selectedFile := c.QueryParam("file")

	if selectedFile == "" {
		selectedFile = fmt.Sprintf("gorgon-%s.log", time.Now().In(time.FixedZone("BRT", -3*60*60)).Format("2006-01-02"))
	}

	entries := readLogFile(filepath.Join(logsPath, selectedFile))
	entries = filterLogs(entries, search, levelFilter)

	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	var files []LogFileInfo
	entriesDir, err := os.ReadDir(logsPath)
	if err == nil {
		for _, f := range entriesDir {
			if f.IsDir() {
				continue
			}
			name := f.Name()
			if !strings.HasPrefix(name, "gorgon-") {
				continue
			}
			info, _ := f.Info()
			size := info.Size()
			sizeStr := formatSize(size)
			isCurrent := name == fmt.Sprintf("gorgon-%s.log", time.Now().In(time.FixedZone("BRT", -3*60*60)).Format("2006-01-02"))
			if strings.HasSuffix(name, ".gz") {
				isCurrent = false
			}
			files = append(files, LogFileInfo{
				Name:      name,
				Size:      sizeStr,
				IsCurrent: isCurrent,
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name > files[j].Name
	})

	return views.Render(c, views.View{
		Layout:    "layout",
		Default:   "logs",
		Templates: map[string]string{"log-entries": "log-entries"},
		Data: LogsData{
			Files:       files,
			Entries:     entries,
			CurrentFile: selectedFile,
			Search:      search,
			Level:       levelFilter,
		},
		Styles: []string{"logs.css"},
	})
}

func readLogFile(path string) []LogEntry {
	f, err := os.Open(path)
	if err != nil {
		gzPath := path + ".gz"
		f, err = os.Open(gzPath)
		if err != nil {
			return nil
		}
		defer f.Close()
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return nil
		}
		defer gzr.Close()
		return parseLogLines(gzr)
	}
	defer f.Close()
	return parseLogLines(f)
}

func parseLogLines(r io.Reader) []LogEntry {
	var entries []LogEntry
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		var raw map[string]any
		if err := json.Unmarshal(line, &raw); err != nil {
			continue
		}

		entry := LogEntry{}

		if t, ok := raw["time"].(string); ok {
			parsed, err := time.Parse(time.RFC3339Nano, t)
			if err == nil {
				entry.Time = parsed.Format("2006-01-02 15:04:05")
			} else {
				entry.Time = t
			}
		}
		if l, ok := raw["level"].(string); ok {
			entry.Level = l
		}
		if m, ok := raw["msg"].(string); ok {
			entry.Message = m
		}

		delete(raw, "time")
		delete(raw, "level")
		delete(raw, "msg")

		if len(raw) > 0 {
			b, _ := json.Marshal(raw)
			entry.Context = string(b)
		}

		entries = append(entries, entry)
	}
	return entries
}

func filterLogs(entries []LogEntry, search, level string) []LogEntry {
	if search == "" && level == "" {
		return entries
	}
	var filtered []LogEntry
	for _, e := range entries {
		if level != "" && !strings.EqualFold(e.Level, level) {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(e.Message), strings.ToLower(search)) && !strings.Contains(e.Context, search) {
			continue
		}
		filtered = append(filtered, e)
	}
	return filtered
}

func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
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
