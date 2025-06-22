package dtos

import (
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	seasonModel "github.com/jusoaresg/gorgon/internal/season/model"
	showModel "github.com/jusoaresg/gorgon/internal/show/model"
	"strconv"
)

type ShowDto struct {
	TvMazeID  int64  `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Language  string `json:"language"`
	Status    string `json:"status"`
	Premiered string `json:"premiered"`
	Ended     string `json:"ended"`
	Rating    struct {
		Average *float64 `json:"average"`
	} `json:"rating"`
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
	Url          string `json:"url"`
	Name         string `json:"name,omitempty"`
	Number       int    `json:"number"`
	PremiereDate string `json:"premiereDate"`
	EndDate      string `json:"endDate"`
}

func (s *SeasonDto) ToModel(showID int64) *seasonModel.Season {
	return &seasonModel.Season{
		ShowID: showID,
		Number: s.Number,
	}
}

func SeasonDtoSliceToModel(seasonDto []SeasonDto, showID int64) []seasonModel.Season {
	var seasons []seasonModel.Season
	for _, season := range seasonDto {
		seasons = append(seasons, *season.ToModel(showID))
	}
	return seasons
}

// Always use as slice
type EpisodeDto struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	Season   int    `json:"season"`
	Number   int    `json:"number"`
	AirStamp string `json:"airStamp"`
	Summary  string `json:"summary"`
	Tracking string
}

func (s *EpisodeDto) ToModel(showID int64, seasonID int64) *episodeModel.Episode {
	return &episodeModel.Episode{
		ShowID:   showID,
		Name:     s.Name,
		SeasonID: seasonID,
		Season:   s.Season,
		Number:   s.Number,
		AirStamp: s.AirStamp,
		Summary:  s.Summary,
		Tracking: s.Tracking,
	}
}

func EpisodesDtoSliceToModel(episodeDto []EpisodeDto, showID int64, seasonID int64) []episodeModel.Episode {
	var episodes []episodeModel.Episode
	for _, episode := range episodeDto {
		episodes = append(episodes, *episode.ToModel(showID, seasonID))
	}
	return episodes
}

func (d *ShowDto) CreateDto(
	TvMazeID int64, Updated int,
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
		TvMazeID:  TvMazeID,
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

func (d *ShowDto) ToModel() showModel.Show {
	imdb, _ := strconv.Atoi(d.Externals.Imdb)

	show := showModel.Show{
		TvMazeID:  d.TvMazeID,
		Name:      d.Name,
		Type:      d.Type,
		Language:  d.Language,
		Status:    d.Status,
		Premiered: d.Premiered,
		Ended:     d.Ended,
		Rating:    d.Rating.Average,
		Summary:   d.Summary,
		Updated:   d.Updated,

		TvRage:   d.Externals.TvRage,
		TheTvDBD: d.Externals.TheTvdb,
		Imdb:     imdb,

		ImageOriginal: d.Image.Original,
		ImageMedium:   d.Image.Medium,
	}
	return show
}
