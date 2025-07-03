package main

import (
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/routes"
	"github.com/jusoaresg/gorgon/internal/scheduler"
	"github.com/jusoaresg/gorgon/internal/scheduler/cron"
	"github.com/jusoaresg/gorgon/internal/scheduler/workers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title           Gongon
// @version         0.1
// @description     Show Download Manager API
// @BasePath /api/v1

// @contact.name   Jusoares
// @contact.email  julianosgreg@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	if err := config.Init(); err != nil {
		panic(fmt.Errorf("failed to initialize config: %w", err))
	}

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

	routes.SetupFrontRouter(e)

	routes.InitializeRoutes(e)
	cron.StartDailyUpdate(scheduler.UpdateAllShows)

	go workers.StartEpisodeSyncWorker(2)
	go workers.VerifySnatchedDownloadsWorker(5)

	cron.StartVerifyEpisodeWasDeleted(scheduler.VerifyEpisodeWasDeleted)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Port)))
}
