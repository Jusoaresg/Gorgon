package web

import (
	"net/http"

	"github.com/jusoaresg/gorgon/internal/show"
	"github.com/labstack/echo/v4"
)

func RegisterShowRoutes(e *echo.Echo, deps *show.Dependencies) {
	handler := NewHandler(deps)

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/shows")
	})

	// Pages
	e.GET("/shows", handler.ShowsRoute)
	e.GET("/show/:id", handler.ShowRoute)
	e.GET("/add-show", handler.AddShowRoute)
	e.GET("/add-show/:tvmaze-id/config", handler.AddShowConfigRoute)
	e.GET("/calendar", handler.CalendarRoute)
	e.GET("/logs", handler.LogsRoute)
	e.GET("/settings", handler.SettingsRoute)
	e.GET("/settings/:type", handler.SettingsTypeRoute)

	// HTMX Partials
	front := e.Group("/front/")
	front.POST("search-show", handler.SearchShowHTMX)
	front.POST("add-show", handler.AddShowHTMX)
	front.GET("show/:id/modal/edit", handler.EditShowModalHTMX)
	front.GET("season/:id/modal/tracking", handler.ChangeSeasonTrackingModal)
	front.GET("calendar/week", handler.CalendarWeekHTMX)
}
