package views

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, view View) error {
	template := view.Default

	if name := c.QueryParam("view"); name != "" {
		if t, ok := view.Templates[name]; ok {
			template = t
		}
	}

	if template != view.Default {
		return c.Render(http.StatusOK, template, PageData{
			Data: view.Data,
		})
	}

	return c.Render(http.StatusOK, view.Layout, PageData{
		TemplateName: template,
		Data:         view.Data,
		Styles:       view.Styles,
	})
}
