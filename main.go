package main

import (
	"fmt"
	"gorgon/config"
	"gorgon/pkg/routes"

	"github.com/gin-gonic/gin"
)

// @title           Gongon API
// @version         0.1
// @description     This is a sample server celler server.

// @contact.name   Jusoares
// @contact.email  julianosgreg@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	g := gin.Default()

	config.Init()

	routes.InitializeRoutes(g)

	g.Run(fmt.Sprintf(":%s", config.Port))
}
