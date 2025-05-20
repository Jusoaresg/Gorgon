package model

import (
	"time"
)

type TrackingStatus string

const (
	StatusWanted     TrackingStatus = "wanted"
	StatusMissing    TrackingStatus = "missing"
	StatusSkipped    TrackingStatus = "skipped"
	StatusSnatched   TrackingStatus = "snatched"
	StatusDownloaded TrackingStatus = "downloaded"
)

type trackingStatusEnum struct{}

var Tracking = trackingStatusEnum{}

func (trackingStatusEnum) Wanted() TrackingStatus     { return StatusWanted }
func (trackingStatusEnum) Missing() TrackingStatus    { return StatusMissing }
func (trackingStatusEnum) Skipped() TrackingStatus    { return StatusSkipped }
func (trackingStatusEnum) Snatched() TrackingStatus   { return StatusSnatched }
func (trackingStatusEnum) Downloaded() TrackingStatus { return StatusDownloaded }

type Episode struct {
	ID       uint `gorm:"primaryKey"`
	ShowId   int  `gorm:"index"`
	Name     string
	Summary  string
	Type     string
	Number   int
	Season   int
	AirStamp string

	FilePath    *string
	Tracking    TrackingStatus `gorm:"type:text;default:'wanted'"`
	TorrentHash string
	Content     []EpisodeContent `gorm:"foreignKey:EpisodeId"`
}

type EpisodeContent struct {
	ID        uint `gorm:"primaryKey"`
	EpisodeId int  `gorm:"index"`

	Name    string  `json:"name"`
	Size    float64 `json:"size"`
	Is_Seed bool    `json:"is_seed"`
}

func (e *Episode) Create(
	ShowId int,
	Name string,
	Summary string,
	Type string,
	Number int,
	Season int,
	AirStamp string,
) *Episode {
	return &Episode{
		ShowId:   ShowId,
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
	e.Tracking = Tracking.Skipped()
	e.FilePath = nil
	e.TorrentHash = ""
}
