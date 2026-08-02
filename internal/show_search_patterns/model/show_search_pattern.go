package model

// ShowSearchPattern is a per-show search query template that is combined
// with the selected filter profile's search patterns on Prowlarr.
type ShowSearchPattern struct {
	ID       int64  `db:"id"`
	ShowID   int64  `db:"show_id"`
	Pattern  string `db:"pattern"`
	Position int    `db:"position"`
}
