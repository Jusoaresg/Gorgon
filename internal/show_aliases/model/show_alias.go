package model

type ShowAlias struct {
	ID      int64  `db:"id"`
	ShowID  int64  `db:"show_id"`
	Alias   string `db:"alias"`
	Country string `db:"country"`
	Source  string `db:"source"`
}
