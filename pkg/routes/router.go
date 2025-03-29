package routes

import (
	"gorgon/docs"
	"gorgon/pkg/handler"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "teste",
	})
}

func InitializeRoutes(g *gin.Engine) {
	basePath := "/api/v1"

	handler.InitHandler()

	docs.SwaggerInfo.BasePath = basePath

	v1 := g.Group(basePath)
	{
		v1.GET("test", test)
	}

	SetupAnilistRouter(v1)

	SetupDatabaseRouter(v1)

	SetupProwlarrRouter(v1)

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}
