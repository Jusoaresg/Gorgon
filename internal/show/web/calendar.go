package web

type CalendarEpisode struct {
	ID         int64  `db:"id"`
	ShowID     int64  `db:"show_id"`
	ShowName   string `db:"show_name"`
	ShowImage  string `db:"show_image"`
	ShowStatus string `db:"show_status"`
	Name       string `db:"name"`
	Number     int    `db:"number"`
	Season     int    `db:"season"`
	AirStamp   int64  `db:"airstamp"`
	Tracking   string `db:"tracking"`
}

type CalendarDay struct {
	Date        string
	DayName     string
	DisplayDate string
	Episodes    []CalendarEpisode
}

type CalendarData struct {
	Days          []CalendarDay
	WeekStart     string
	WeekEnd       string
	PrevWeek      string
	NextWeek      string
	Today         string
	IsCurrentWeek bool
}
