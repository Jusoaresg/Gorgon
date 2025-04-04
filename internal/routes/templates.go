package routes

import (
	"gorgon/config"
	"gorgon/internal/templates/handler"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func SetupTemplatesRouter(r *echo.Group) {
	logger := config.GetLogger()

	r.Static("/css", "./assets/css")

	templatesRouter := r.Group("")
	{
		templatesRouter.GET("/", handler.RenderShowsList)
		logger.Info("GET route added to /")
		templatesRouter.POST("/", handler.RenderShowsList)
		logger.Info("POST route added to /")

		templatesRouter.GET("/add-show", handler.RenderAddShow)
		logger.Info("GET route added to /add-show")
		templatesRouter.POST("/add-show", handler.RenderAddShow)
		logger.Info("POST route added to /add-show")

		templatesRouter.GET("/show", handler.RenderShow)
		logger.Info("GET route added to /show")

	}
	logger.Info("Templates routes added successfully", slog.String("endpoint", "/"))
}
