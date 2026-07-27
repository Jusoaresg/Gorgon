package routes

import (
	"log/slog"
	"net/http"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/docs"
	"github.com/jusoaresg/gorgon/internal/app"
	showRouter "github.com/jusoaresg/gorgon/internal/show/web"
	"github.com/jusoaresg/gorgon/pkg/handler"
	"github.com/labstack/echo/v4"
)

const basePath = "/api/v1/"

func InitializeRoutes(e *echo.Echo, deps *app.Dependencies) {
	logger := config.GetLogger().WithGroup("routes").With("name", "initializeRoutes")

	handler.InitHandler()

	docs.SwaggerInfo.BasePath = basePath
	logger.Info("Initializing routes", slog.String("basePath", basePath))

	showRouter.RegisterShowRoutes(e, deps.Show)
	logger.Debug("Show web routes initialized successfully")

	v1 := e.Group(basePath)

	SetupTvMazeRouter(v1)
	logger.Debug("TvMaze route initialized successfully")

	SetupDatabaseRouter(v1, deps.Show, deps.Episode, deps.Season, deps.Indexer, deps.ShowAliases)
	logger.Debug("Database route initialized successfully")

	SetupProwlarrRouter(v1)
	logger.Debug("Prowlarr route initialized successfully")

	SetupQbittorrentRouter(v1)
	logger.Debug("QBittorrent route initialized successfully")

	SetupWebsocketRouter(v1)
	logger.Debug("WebSocket route initialized successfully")

	SetupAppConfigRouter(v1, deps.AppConfig)
	logger.Debug("App Config route initialized successfully")

	e.GET("/swagger.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", docs.SwaggerJSON)
	})

	v1.GET("docs", func(c echo.Context) error {
		return c.Render(http.StatusOK, "docs.html", nil)
	})

	logger.Info("Routes initialized", slog.String("basePath", basePath))
}
