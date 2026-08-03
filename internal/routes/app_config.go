package routes

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/config"
	appConfig "github.com/jusoaresg/gorgon/internal/config"
	appConfigApi "github.com/jusoaresg/gorgon/internal/config/api"

	"github.com/labstack/echo/v4"
)

func SetupAppConfigRouter(v1 *echo.Group, deps *appConfig.Dependencies) {
	logger := config.GetLogger()

	handler := appConfigApi.NewHandler(deps)

	appConfigGroup := v1.Group("app/config")
	{
		appConfigGroup.Match([]string{"POST", "PATCH"}, "", handler.UpdateAppConfig)
		appConfigGroup.GET("", handler.GetAppConfig)
	}
	logger.Info("App Config routes added successfully", slog.String("endpoint", "/api/v1/app/config"))
}
