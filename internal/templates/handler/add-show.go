package handler

import (
	"fmt"
	"gorgon/assets/templates/components"
	"gorgon/assets/templates/pages"
	"gorgon/config"
	"gorgon/external/tvmaze/schema"
	"gorgon/external/tvmaze/service"

	"github.com/labstack/echo/v4"
)

func AddHandler(c echo.Context) error {
	logger := config.GetLogger()

	if c.Request().Method == echo.POST {
		query := c.FormValue("query")
		fmt.Println(query)

		tvMazeService := service.NewTvMazeSearchService(logger)

		var shows []schema.TvMazeResponse
		tvMazeService.SearchByName(query, &shows)

		component := components.SearchResults(shows)

		component.Render(c.Request().Context(), c.Response())
		return nil
	}

	page := pages.AddShow("Add Show", []schema.TvMazeResponse{})
	page.Render(c.Request().Context(), c.Response())

	return nil
}
