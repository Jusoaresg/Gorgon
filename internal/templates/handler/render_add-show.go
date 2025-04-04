package handler

import (
	"fmt"
	"gorgon/assets/templates/components/search"
	"gorgon/assets/templates/pages"
	"gorgon/config"
	"gorgon/external/tvmaze/schema"
	"gorgon/external/tvmaze/service"
	"gorgon/internal/db/model"
	"gorgon/pkg/services"

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

		var ids []int
		for _, show := range *shows {
			ids = append(ids, show.Show.ShowID)
		}

		var addedShows []model.Show
		baseService := services.NewBaseService()
		if err := baseService.GetShowsByIdentification(&addedShows, "show_id", ids); err != nil {
			return err
		}

		addedMap := make(map[int]bool)
		for _, show := range addedShows {
			addedMap[show.ShowID] = true
		}

		component := search.SearchResults(*shows, addedMap)
		component.Render(c.Request().Context(), c.Response())
		return nil
	}

	page := pages.AddShow("Add Show", []schema.TvMazeResponse{})
	page.Render(c.Request().Context(), c.Response())

	return nil
}
