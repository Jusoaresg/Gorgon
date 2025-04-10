package handler

import (
	"fmt"
	"gorgon/assets/templates/components/shows_list"
	"gorgon/assets/templates/pages"
	"gorgon/internal/db/model"
	"gorgon/internal/templates"
	"gorgon/pkg/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RenderShowsList(c echo.Context) error {
	baseService := services.NewBaseService()

	var shows []model.Show
	if c.Request().Method == echo.POST {
		name := c.FormValue("query")
		fmt.Println(name)
		if err := baseService.ListByNameWithPreload(name, &shows, "Seasons", "Episodes"); err != nil {
			return err
		}

		fmt.Println(shows)
		component := shows_list.ShowGridShowsList(shows)
		component.Render(c.Request().Context(), c.Response())
		return nil

	} else {
		if err := baseService.ListWithPreload(&shows, "Seasons", "Episodes"); err != nil {
			return err
		}
	}

	page := pages.ShowsList("Shows List", shows)
	if err := templates.Render(c, page); err != nil {
		return err
	}

	return nil
}

func RedirectToShow(c echo.Context) error {
	id := c.FormValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	service := services.NewBaseService()

	var show model.Show
	if err := service.GetWithPreload(&show, idInt, "Seasons", "Episodes"); err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/show/%d", show.ID))
	return nil
}
