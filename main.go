package main

import (
	"fmt"
	"gorgon/config"
	"gorgon/internal/db/model"
	"gorgon/internal/routes"
	"gorgon/internal/scheduler"
	"gorgon/internal/scheduler/cron"
	"gorgon/pkg/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title           Gongon
// @version         0.1
// @description     Anime download manager API
// @BasePath /api/v1

// @contact.name   Jusoares
// @contact.email  julianosgreg@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	config.Init()

	e := echo.New()

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "./assets/front/build/",
		Browse:     true,
		HTML5:      true,
		IgnoreBase: true,
	}))

	cors := middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	e.Use(middleware.CORSWithConfig(cors))
	e.Use(middleware.Logger())

	e.Static("/static", "static")

	e.GET("/*", func(c echo.Context) error {
		return c.File("./assets/front/build/index.html")
	})

	e.GET("/", func(c echo.Context) error {
		baseService := services.NewBaseService()

		var shows []model.Show
		baseService.List(&shows)

		return c.JSON(200, shows)
	})

	routes.InitializeRoutes(e)
	cron.StartDailyUpdate(scheduler.UpdateAllShows)
	cron.StartSearchNewEpisodes(scheduler.SyncWantedEpisodes)
	cron.StartVerifySnatched(scheduler.VerifySnatchedDownload)
	cron.StartVerifyEpisodeWasDeleted(scheduler.VerifyEpisodeWasDeleted)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Port)))
}
