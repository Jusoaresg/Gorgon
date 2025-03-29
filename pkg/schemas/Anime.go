package schemas

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Anime struct {
	gorm.Model
	Aid          string `json:"id"`
	EpisodeCount int    `json:"episodes"`
	Description  string
	AiringTime   string `json:"nextairingepisode"`

	InstalledEps datatypes.JSON `gorm:"type:json"`

	Titles Title          `gorm:"foreignKey:Aid"`
	Genres datatypes.JSON `json:"genres"`
}

type Title struct {
	gorm.Model
	Aid     string
	English string
	Romaji  string
}
