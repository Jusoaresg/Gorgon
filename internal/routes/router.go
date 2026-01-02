package routes

import (
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/docs"
	"github.com/jusoaresg/gorgon/pkg/handler"
	"github.com/labstack/echo/v4"
)

const basePath = "/api/v1/"

func InitializeRoutes(e *echo.Echo) {
	logger := config.GetLogger().WithGroup("routes").With("name", "initializeRoutes")

	handler.InitHandler()

	docs.SwaggerInfo.BasePath = basePath

	logger.Info("Initializing routes", slog.String("basePath", basePath))

	v1 := e.Group(basePath)

	SetupTvMazeRouter(v1)
	logger.Debug("TvMaze route initialized successfully")

	SetupDatabaseRouter(v1)
	logger.Debug("Database route initialized successfully")

	SetupProwlarrRouter(v1)
	logger.Debug("Prowlarr route initialized successfully")

	SetupQbittorrentRouter(v1)
	logger.Debug("QBittorrent route initialized successfully")

	SetupWebsocketRouter(v1)
	logger.Debug("WebSocket route initialized successfully")

	SetupAppConfigRouter(v1)
	logger.Debug("App Config route initialized successfully")

	logger.Info("Routes initialized", slog.String("basePath", basePath))

	e.GET("/swagger.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", docs.SwaggerJSON)
	})

	v1.GET("docs", func(c echo.Context) error {
		return c.Render(http.StatusOK, "docs.html", nil)
	})
}
