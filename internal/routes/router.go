package routes

import (
	_ "embed"
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/docs"
	tvMazeService "github.com/jusoaresg/gorgon/external/tvmaze/service"
	showAggregator "github.com/jusoaresg/gorgon/internal/show/service"
	"github.com/jusoaresg/gorgon/pkg/handler"
	"github.com/labstack/echo/v4"
)

type RoutersDeps struct {
	Db             *sqlx.DB
	Logger         *slog.Logger
	AggShowService *showAggregator.ShowAggregatorService
	TvMazeService  *tvMazeService.TvMazeSearchService
}

const basePath = "/api/v1/"

func InitializeRoutes(e *echo.Echo, deps *RoutersDeps) {
	logger := config.GetLogger().WithGroup("routes").With("name", "initializeRoutes")

	handler.InitHandler()

	docs.SwaggerInfo.BasePath = basePath
	logger.Info("Initializing routes", slog.String("basePath", basePath))

	SetupFrontRouter(e)
	logger.Debug("Front route initialized successfully")

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

	e.GET("/swagger.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", docs.SwaggerJSON)
	})

	v1.GET("docs", func(c echo.Context) error {
		return c.Render(http.StatusOK, "docs.html", nil)
	})

	logger.Info("Routes initialized", slog.String("basePath", basePath))
}
