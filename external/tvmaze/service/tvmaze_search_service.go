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

func (t *TvMazeSearchService) SearchByName(name string) (*[]schema.TvMazeResponse, error) {
	var model []schema.TvMazeResponse
	if err := t.APIService.Get(fmt.Sprintf("/search/shows?q=%s", url.QueryEscape(name)), &model); err != nil {
		t.Logger.Error("Error while searching for name", slog.String("error", err.Error()))
		return nil, err
	}
	return &model, nil
}

func (t *TvMazeSearchService) SearchByTvMazeId(id int) (*dtos.ShowDto, error) {
	var model dtos.ShowDto
	if err := t.APIService.Get(fmt.Sprintf("/shows/%d", id), &model); err != nil {
		return nil, err
	}
	return &model, nil
}

func (t *TvMazeSearchService) SearchByTheTvDbId(id string) (*dtos.ShowDto, error) {
	var model dtos.ShowDto
	if err := t.APIService.Get(fmt.Sprintf("/lookup/shows?thetvdb=%s", id), &model); err != nil {
		return nil, err
	}
	return &model, nil
}

func (t *TvMazeSearchService) SearchByImdb(id string) (*dtos.ShowDto, error) {
	var model dtos.ShowDto
	if err := t.APIService.Get(fmt.Sprintf("/lookup/shows?imdb=%s", id), &model); err != nil {
		return nil, err
	}
	return &model, nil
}

// Episodes
func (t *TvMazeSearchService) SearchEpisodes(id int) (*[]dtos.EpisodeDto, error) {
	var model []dtos.EpisodeDto
	if err := t.APIService.Get(fmt.Sprintf("/shows/%d/episodes", id), &model); err != nil {
		return nil, err
	}
	return &model, nil
}

// Seasons
func (t *TvMazeSearchService) SearchSeasons(id int) (*[]dtos.SeasonDto, error) {
	var model []dtos.SeasonDto
	if err := t.APIService.Get(fmt.Sprintf("/shows/%d/seasons", id), &model); err != nil {
		return nil, err
	}
	return &model, nil
}
