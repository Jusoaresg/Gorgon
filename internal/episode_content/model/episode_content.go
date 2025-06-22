package model

type EpisodeContent struct {
	ID        int64 `db:"id"`
	EpisodeId int64 `db:"episode_id"`

	Name     string  `db:"name"`
	FilePath string  `db:"file_path"`
	Size     float64 `db:"size"`
	Is_Seed  bool    `db:"is_seed"`
}
