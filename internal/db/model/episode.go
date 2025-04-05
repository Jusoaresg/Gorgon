package model

import "time"

type TrackingStatus string

const (
	StatusWanted     TrackingStatus = "wanted"
	StatusMissing    TrackingStatus = "missing"
	StatusSkipped    TrackingStatus = "skipped"
	StatusSnatched   TrackingStatus = "snatched"
	StatusDownloaded TrackingStatus = "downloaded"
)

type Episode struct {
	ID       uint `gorm:"primaryKey"`
	ShowId   int  `gorm:"index"`
	Name     string
	Summary  string
	Type     string
	Number   int
	Season   int
	AirStamp string

	Installed bool
	FilePath  *string
	Tracking  TrackingStatus `gorm:"type:text;default:'wanted'"`
}

func (e *Episode) AirDate() string {
	t, _ := time.Parse(time.RFC3339, e.AirStamp)
	return t.Format("2006-01-02")
}

func (e *Episode) AirTime() string {
	t, _ := time.Parse(time.RFC3339, e.AirStamp)
	return t.Format("15:04")
}
