package testutils

import (
	"fmt"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	seasonModel "github.com/jusoaresg/gorgon/internal/season/model"
	showModel "github.com/jusoaresg/gorgon/internal/show/model"
	"math/rand"
	"time"
)

func MakeFakeEpisode() episodeModel.Episode {
	idSuffix := rand.Intn(10000)

	return episodeModel.Episode{
		ShowID:   rand.Int63n(1_000_000),
		SeasonID: rand.Int63n(1_000_000),
		Name:     fmt.Sprintf("Test Episode %d", idSuffix),
		Summary:  fmt.Sprintf("Summary for test episode %d", idSuffix),
		Type:     "Episode Type",
		Number:   rand.Int(),
		Season:   rand.Int(),
		AirStamp: time.Now().GoString(),
	}
}

func MakeFakeShow() showModel.Show {
	rating := rand.Float64()*10 + 0.1
	idSuffix := rand.Intn(10000)

	return showModel.Show{
		TvMazeID:      rand.Int63n(1_000_000),
		Name:          fmt.Sprintf("Test Show %d", idSuffix),
		Type:          "Series",
		Language:      "English",
		Status:        "Running",
		Premiered:     "2023-01-01",
		Ended:         "",
		Rating:        &rating,
		Summary:       fmt.Sprintf("Summary for test show %d", idSuffix),
		Updated:       int(time.Now().Unix()),
		TvRage:        rand.Int(),
		TheTvDBD:      rand.Int(),
		Imdb:          rand.Int(),
		ImageOriginal: "https://img/original.jpg",
		ImageMedium:   "https://img/medium.jpg",
		Genres:        "Action,Drama",
	}
}

func MakeFakeSeason() seasonModel.Season {
	return seasonModel.Season{
		Number: rand.Int(),
	}
}
