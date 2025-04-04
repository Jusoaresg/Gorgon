package handler

import (
	"gorgon/assets/templates/pages"
	"gorgon/internal/db/model"
	"gorgon/internal/templates"
	"gorgon/pkg/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RenderShow(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	service := services.NewBaseService()
	var show model.Show
	if err := service.GetWithPreload(&show, idInt, "Seasons", "Episodes"); err != nil {
		return err
	}

	page := pages.Show("Show", show)
	// return page.Render(c.Request().Context(), c.Response())
	return templates.Render(c, page)
}
