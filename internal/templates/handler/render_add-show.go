package handler

import (
	"fmt"
	"gorgon/assets/templates/components/search"
	"gorgon/assets/templates/pages"
	"gorgon/config"
	"gorgon/external/tvmaze/schema"
	"gorgon/external/tvmaze/service"

	"github.com/labstack/echo/v4"
)

func RenderAddShow(c echo.Context) error {
	logger := config.GetLogger()

	if c.Request().Method == echo.POST {
		query := c.FormValue("query")
		fmt.Println(query)

		tvMazeService := service.NewTvMazeSearchService(logger)

		shows, err := tvMazeService.SearchByName(query)
		if err != nil {
			return err
		}

		component := search.SearchResults(*shows)

		component.Render(c.Request().Context(), c.Response())
		return nil
	}

	page := pages.AddShow("Add Show", []schema.TvMazeResponse{})
	page.Render(c.Request().Context(), c.Response())

	return nil
}
