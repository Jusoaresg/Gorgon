package dtos

import (
	"gorgon/internal/db/model"
)

type ShowDto struct {
	ShowID    int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Language  string `json:"language"`
	Status    string `json:"status"`
	Premiered string `json:"premiered"`
	Ended     string `json:"ended"`
	Rating    struct {
		Average *float64 `json:"average"`
	} `json:"rating,omitempty"`
	Summary string `json:"summary"`
	Updated int    `json:"updated"`

	Schedule struct {
		Days []string `json:"days"`
		Time string   `json:"time"`
	} `json:"schedule"`
	Externals struct {
		Imdb    string `json:"imdb"`
		TheTvdb int    `json:"thetvdb"`
		TvRage  int    `json:"tvrage"`
	} `json:"externals"`
	Image struct {
		Medium   string `json:"medium"`
		Original string `json:"original"`
	} `json:"image"`
	Genres []string `json:"genres"`
}

// Always use as slice
type SeasonDto struct {
	ShowId       int    `json:"id"`
	Url          string `json:"url"`
	Name         string `json:"name,omitempty"`
	Number       int    `json:"number"`
	PremiereDate string `json:"premiereDate"`
	EndDate      string `json:"endDate"`
}

// Always use as slice
type EpisodeDto struct {
	ShowId   int    `json:"id"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	Season   int    `json:"season"`
	Number   int    `json:"number"`
	AirStamp string `json:"airStamp"`
	Summary  string `json:"summary"`
	Tracking model.TrackingStatus
}

func (d *ShowDto) CreateDto(
	ShowId, Updated int,
	Name, Type, Language, Status, Premiered, Ended, Summary string,
	Rating struct {
		Average *float64 `json:"average"`
	},
	Scedule struct {
		Days []string
		Time string
	},
	Externals struct {
		Imdb    string `json:"imdb"`
		TheTvdb int    `json:"thetvdb"`
		TvRage  int    `json:"tvrage"`
	},
	Image struct {
		Medium   string `json:"medium"`
		Original string `json:"original"`
	},
	Genres []string,
) *ShowDto {

	return &ShowDto{
		ShowID:    ShowId,
		Name:      Name,
		Type:      Type,
		Language:  Language,
		Status:    Status,
		Premiered: Premiered,
		Ended:     Ended,
		Summary:   Summary,
		Rating:    Rating,
		Externals: Externals,
		Image:     Image,
		Genres:    Genres,
	}
}

func (d *ShowDto) ToModel(episodes *[]EpisodeDto, seasons *[]SeasonDto) *model.Show {

	show := model.Show{
		ShowID:    d.ShowID,
		Name:      d.Name,
		Type:      d.Type,
		Language:  d.Language,
		Status:    d.Status,
		Premiered: d.Premiered,
		Ended:     d.Ended,
		Rating:    d.Rating.Average,
		Summary:   d.Summary,
		Updated:   d.Updated,

		Seasons:  make([]model.Season, len(*seasons)),
		Episodes: make([]model.Episode, len(*episodes)),

		Externals: model.Externals{
			Tvrage:   d.Externals.TvRage,
			Thetvdvb: d.Externals.TheTvdb,
			Imdb:     d.Externals.Imdb,
		},
		Image: model.Image{
			Original: d.Image.Original,
			Medium:   d.Image.Medium,
		},
	}

	for i, season := range *seasons {
		show.Seasons[i] = model.Season{
			ShowId:   d.ShowID,
			SeasonId: season.ShowId,
			Number:   season.Number,
		}
	}

	for i, episode := range *episodes {
		show.Episodes[i] = model.Episode{
			ShowId:   episode.ShowId,
			Name:     episode.Name,
			Summary:  episode.Summary,
			Number:   episode.Number,
			Season:   episode.Season,
			AirStamp: episode.AirStamp,
			Tracking: episode.Tracking,
		}
	}

	return &show
}
