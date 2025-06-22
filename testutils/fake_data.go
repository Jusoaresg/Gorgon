package testutils

import (
	"fmt"
	"github.com/jusoaresg/gorgon/internal/db/model"
	"math/rand"
	"time"
)

func MakeFakeEpisode() model.Episode {
	idSuffix := rand.Intn(10000)

	return model.Episode{
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

func MakeFakeShow() model.Show {
	rating := rand.Float64()*10 + 0.1
	idSuffix := rand.Intn(10000)

	return model.Show{
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

func MakeFakeSeason() model.Season {
	return model.Season{
		Number: rand.Int(),
	}
}
