package model

type Show struct {
	ID        int64    `db:"id"`
	TvMazeID  int64    `db:"tv_maze_id"`
	Name      string   `db:"name"`
	Type      string   `db:"type"`
	Language  string   `db:"language"`
	Status    string   `db:"status"`
	Premiered string   `db:"premiered"`
	Ended     string   `db:"ended"`
	Rating    *float64 `db:"rating"`
	Summary   string   `db:"summary"`
	Updated   int      `db:"updated"`

	TvRage   int `db:"tv_rage"`
	TheTvDBD int `db:"the_tvdbd"`
	Imdb     int `db:"imdb"`

	ImageOriginal string `db:"image_original"`
	ImageMedium   string `db:"image_medium"`

	Genres string `db:"genres"`
}

type Schedule struct {
	ID     int64  `db:"id"`
	ShowId int64  `db:"show_id"`
	Time   string `db:"time"`
	Days   string `db:"days"`
}
