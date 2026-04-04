package handler

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Update App Config
// @Description Update Gorgon Application Config
// @Tags App/Config
// @Produce json
// @Param request body schemas.ConfigFile true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /app/config [post]
// @Router /app/config [patch]
func UpdateAppConfig(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Update App Config", slog.String("endpoint", "/app/config"), slog.String("method", c.Request().Method))

	var request schemas.UpdateConfigInput

	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 500, "Failed to bind body")
		return err
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load app config file")
		schemas.SendError(c, 500, "Failed to load app config file")
		return err
	}

	cfg.Apply(&request)

	if err := config.SaveConfig(cfg); err != nil {
		logger.Error("Failed to save app config file")
		schemas.SendError(c, 500, "Failed to save app config file")
		return err
	}

	logger.Info("Successfully updated app config")
	schemas.SendSuccess(c, "Update App Config", cfg)
	return nil
}
