package handler

import (
	"fmt"
	"gorgon/assets/templates/pages"
	"gorgon/internal/db/model"
	"gorgon/pkg/services"

	"github.com/labstack/echo/v4"
)

func RenderShow(c echo.Context) error {

	service := services.NewBaseService()
	var show model.Show
	if err := service.GetWithPreload(&show, 7, "Seasons", "Episodes"); err != nil {
		return err
	}
	fmt.Println(show.Seasons)
	page := pages.Show("Show", show)
	page.Render(c.Request().Context(), c.Response())
	return nil
}
