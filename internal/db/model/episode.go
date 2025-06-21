package model

import (
	"time"
)

const (
	TrackingWanted     string = "wanted"
	TrackingMissing    string = "missing"
	TrackingSkipped    string = "skipped"
	TrackingSnatched   string = "snatched"
	TrackingDownloaded string = "downloaded"
)

type Episode struct {
	ID       int64 `db:"id"`
	ShowID   int64 `db:"show_id"`
	SeasonID int64 `db:"season_id"`

	Name     string `db:"name"`
	Summary  string `db:"summary"`
	Type     string `db:"type"`
	Number   int    `db:"number"`
	Season   int    `db:"season"`
	AirStamp string `db:"airstamp"`

	Tracking    string `db:"tracking"`
	TorrentHash string `db:"torrent_hash"`
}

type EpisodeContent struct {
	ID        int   `db:"id"`
	EpisodeId int64 `db:"episode_id"`

	Name     string  `db:"name"`
	FilePath string  `db:"file_path"`
	Size     float64 `db:"size"`
	Is_Seed  bool    `db:"is_seed"`
}

func (e *Episode) Create(
	ShowID int64,
	Name string,
	Summary string,
	Type string,
	Number int,
	Season int,
	AirStamp string,
) *Episode {
	return &Episode{
		ShowID:   ShowID,
		Name:     Name,
		Summary:  Summary,
		Type:     Type,
		Number:   Number,
		Season:   Season,
		AirStamp: AirStamp,
	}
}

func (e *Episode) AirDate() string {
	t, _ := time.Parse(time.RFC3339, e.AirStamp)
	return t.Format("2006-01-02")
}

func (e *Episode) AirTime() string {
	t, _ := time.Parse(time.RFC3339, e.AirStamp)
	return t.Format("15:04")
}

func (e *Episode) HasAired() (bool, error) {
	airTime, err := time.Parse(time.RFC3339, e.AirStamp)
	if err != nil {
		return false, err
	}

	// If the airTime is before or equal the time now, then it has been aired
	return !airTime.After(time.Now()), nil
}

func (e *Episode) SetNotInstalled() {
	e.Tracking = TrackingSkipped
	e.TorrentHash = ""
}
