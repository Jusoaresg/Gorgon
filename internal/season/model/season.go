package model

type Season struct {
	ID     int64 `db:"id"`
	ShowID int64 `db:"show_id"`
	Number int   `db:"season_number"`
}
