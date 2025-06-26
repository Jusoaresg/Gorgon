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
func UpdateAppConfig(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Update App Config", slog.String("endpoint", "/app/config"), slog.String("method", "post"))

	var request schemas.ConfigFile

	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 500, "Failed to bind body")
		return err
	}

	if err := config.SaveConfig(request); err != nil {
		//TODO: Error logger
		schemas.SendError(c, 500, "Failed to save app config file")
		return err
	}

	logger.Info("Successfully updated app config")
	schemas.SendSuccess(c, "Update App Config", request)
	return nil
}
