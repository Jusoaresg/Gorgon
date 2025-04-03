package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"gorgon/config"
	"gorgon/internal/routes"

	"github.com/gin-gonic/gin"
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
	g := gin.Default()

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                // Permitindo apenas a origem do seu frontend (ajuste conforme necessário)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},     // Métodos permitidos
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // Cabeçalhos permitidos
		ExposeHeaders:    []string{"Content-Length"},                   // Cabeçalhos expostos ao frontend
		AllowCredentials: true,                                         // Permitir envio de cookies e credenciais
	}))

	config.Init()

	routes.InitializeRoutes(g)

	g.Run(fmt.Sprintf(":%s", config.Port))
}
