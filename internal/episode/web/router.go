package web

import (
	"github.com/jusoaresg/gorgon/internal/episode"
	"github.com/labstack/echo/v4"
)

func RegisterEpisodeRoutes(e *echo.Echo, deps *episode.Dependencies) {
	handler := NewHandler(deps)

	// HTMX Partials
	front := e.Group("/front/")
	front.GET("episode/:id/modal/tracking", handler.ChangeEpisodeTrackingModal)
	front.GET("episode/:id/interactive-search", handler.InteractiveSearchModal)
	front.GET("episode/:id/search-results", handler.SearchEpisodeResults)
	front.GET("episode/:id/search-alias", handler.SearchAliasResult)
	front.POST("episode/:id/download", handler.DownloadEpisodeTorrent)
}
