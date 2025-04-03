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
		templatesRouter.GET("/", handler.ShowsListHandler)
		logger.Info("GET route added to /")
		templatesRouter.POST("/", handler.ShowsListHandler)
		logger.Info("POST route added to /")

		templatesRouter.GET("/add", handler.AddHandler)
		logger.Info("GET route added to /add")
		templatesRouter.POST("/add", handler.AddHandler)
		logger.Info("POST route added to /add")

	}
	logger.Info("Templates routes added successfully", slog.String("endpoint", "/"))
}
