package routes

import (
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/docs"
	"github.com/jusoaresg/gorgon/pkg/handler"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"log/slog"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
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

	v1.GET("docs", func(c echo.Context) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			CDN:     "https://cdn.jsdelivr.net/npm/@scalar/api-reference@latest",
			SpecURL: "./docs/swagger.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Gorgon",
			},
			DarkMode: true,
		})

		if err != nil {
			logger.Error("Error while generating API documentation", slog.String("error", err.Error()))
			schemas.SendError(c, 500, fmt.Sprintf("Error while getting Scalar Api Reference: %s", err.Error()))
			return err
		}

		c.HTML(200, htmlContent)
		logger.Info("API Documentation generated successfully", slog.String("url", "/docs"))
		return nil
	})
}
