package routes

import (
	"fmt"
	"gorgon/config"
	"gorgon/docs"
	"gorgon/pkg/handler"
	"gorgon/pkg/schemas"
	"log/slog"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(g *gin.Engine) {
	logger := config.GetLogger()

	basePath := "/api/v1"

	handler.InitHandler()

	docs.SwaggerInfo.BasePath = basePath

	logger.Info("Initializing routes", slog.String("basePath", basePath))

	v1 := g.Group(basePath)
	SetupAnilistRouter(v1)
	logger.Debug("Anilist route initialized sucessfully")
	SetupDatabaseRouter(v1)
	logger.Debug("Database route initialized sucessfully")
	SetupProwlarrRouter(v1)
	logger.Debug("Prowlarr route initialized sucessfully")
	SetupQbittorrentRouter(v1)
	logger.Debug("QBittorrent route initialized sucessfully")

	logger.Info("Routes initialized", slog.String("basePath", basePath))

	g.GET("/docs", func(c *gin.Context) {
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
			return
		}

		c.Data(200, "text/html", []byte(htmlContent))
		logger.Info("API Documentation generated successfully", slog.String("url", "/docs"))
	})
}
