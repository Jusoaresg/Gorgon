package routes

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/config/handler"

	"github.com/labstack/echo/v4"
)

func SetupAppConfigRouter(v1 *echo.Group) {
	logger := config.GetLogger()

	appConfigGroup := v1.Group("app/config")
	{
		appConfigGroup.POST("", handler.UpdateAppConfig)
		appConfigGroup.GET("", handler.GetAppConfig)

	}
	logger.Info("App Config routes added successfully", slog.String("endpoint", "/api/v1/app/config"))
}
