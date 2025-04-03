package main

import (
	"encoding/json"
	"fmt"
	"gorgon/config"
	"gorgon/internal/routes"
	"os"

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
	// g := gin.Default()

	e := echo.New()

	cors := middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	e.Use(middleware.CORSWithConfig(cors))

	config.Init()

	routes.InitializeRoutes(e)

	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		return
	}
	os.WriteFile("routes.json", data, 0644)

	// g.Run(fmt.Sprintf(":%s", config.Port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Port)))
}
