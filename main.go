package main

import (
	"embed"
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/routes"
	"github.com/jusoaresg/gorgon/internal/scheduler"
	"github.com/jusoaresg/gorgon/internal/scheduler/cron"
	"github.com/jusoaresg/gorgon/internal/scheduler/workers"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed all:assets/front/build
var embeddedStaticFiles embed.FS

// @title           Gongon
// @version         0.1
// @description     Show Download Manager API
// @BasePath /api/v1

// @contact.name   Jusoares
// @contact.email  julianosgreg@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	config.Init()

	e := echo.New()

	cors := middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	e.Use(middleware.CORSWithConfig(cors))
	e.Use(middleware.Logger())

	buildFS, err := fs.Sub(embeddedStaticFiles, "assets/front/build")
	if err != nil {
		panic(fmt.Errorf("erro ao acessar build embedado: %w", err))
	}

	e.GET("/favicon.png", echo.WrapHandler(http.FileServer(http.FS(buildFS))))
	e.GET("/_app/*", echo.WrapHandler(http.FileServer(http.FS(buildFS))))
	e.GET("/css/*", echo.WrapHandler(http.FileServer(http.FS(buildFS))))

	e.GET("/*", func(c echo.Context) error {
		index, err := fs.ReadFile(buildFS, "index.html")
		if err != nil {
			panic(err)
		}
		return c.Blob(http.StatusOK, echo.MIMETextHTML, index)
	})

	routes.InitializeRoutes(e)
	cron.StartDailyUpdate(scheduler.UpdateAllShows)

	go workers.StartEpisodeSyncWorker(2)
	go workers.VerifySnatchedDownloadsWorker(5)

	cron.StartVerifyEpisodeWasDeleted(scheduler.VerifyEpisodeWasDeleted)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Port)))
}
