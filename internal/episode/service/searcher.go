package service

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/external/prowlarr/service"
)

func SearchEpisode(query string) ([]schema.SearchResponse, error) {
	logger := config.GetLogger()
	prowlarrService, err := service.NewProwlarrSearchService(logger)
	if err != nil {
		return nil, err
	}

	request := schema.SearchRequest{
		Query: query,
	}

	var response []schema.SearchResponse
	err = prowlarrService.Search(&request, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
