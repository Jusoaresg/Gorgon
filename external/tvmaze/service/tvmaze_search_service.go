package service

import (
	"fmt"
	"gorgon/external/tvmaze/schema"
	"gorgon/pkg/schemas/dtos"
	"gorgon/pkg/services"
	"log/slog"
	"net/url"
)

type TvMazeSearchService struct {
	Logger     *slog.Logger
	APIService services.APIService
	Url        string
}

func NewTvMazeSearchService(logger *slog.Logger) *TvMazeSearchService {
	url := "http://api.tvmaze.com"
	return &TvMazeSearchService{
		Logger:     logger,
		Url:        url,
		APIService: *services.NewAPIService(url, logger),
	}
}

func (t *TvMazeSearchService) SearchByName(name string, model *[]schema.TvMazeResponse) error {
	if err := t.APIService.Get(fmt.Sprintf("/search/shows?q=%s", url.QueryEscape(name)), &model); err != nil {
		t.Logger.Error("Error while searching for name", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (t *TvMazeSearchService) SearchByTvMazeId(id int, model *dtos.ShowDto) error {
	if err := t.APIService.Get(fmt.Sprintf("/shows/%d", id), &model); err != nil {
		return err
	}
	return nil
}

func (t *TvMazeSearchService) SearchByTheTvDbId(id string, model *dtos.ShowDto) error {
	if err := t.APIService.Get(fmt.Sprintf("/lookup/shows?thetvdb=%s", t.Url, id), &model); err != nil {
		return err
	}
	return nil
}

func (t *TvMazeSearchService) SearchByImdb(id string, model *dtos.ShowDto) error {
	if err := t.APIService.Get(fmt.Sprintf("/lookup/shows?imdb=%s", t.Url, id), &model); err != nil {
		return err
	}
	return nil
}

// Episodes
func (t *TvMazeSearchService) SearchEpisodes(id int, model *[]dtos.EpisodeDto) error {
	if err := t.APIService.Get(fmt.Sprintf("/shows/%d/episodes", id), &model); err != nil {
		return err
	}
	return nil
}

// Seasons
func (t *TvMazeSearchService) SearchSeasons(id int, model *[]dtos.SeasonDto) error {
	if err := t.APIService.Get(fmt.Sprintf("/shows/%d/seasons", id), &model); err != nil {
		return err
	}
	return nil
}
