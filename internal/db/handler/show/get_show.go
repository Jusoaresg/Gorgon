package show

import (
	"gorgon/config"
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Get Show
// @Description Getshow
// @Tags Database/Show
// @Produce json
// @Param id path int true "Show ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/{id} [get]
func GetShow(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Get Show", slog.String("endpoint", "/database/show/:id"), slog.String("method", "get"))

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		schemas.SendError(c, 400, "Error while converting id to int")
		return err
	}

	var show model.Show
	baseService := services.NewBaseService()

	if err := baseService.GetWithPreload(&show, idInt, "Seasons", "Episodes"); err != nil {
		logger.Error("Error while fetching show from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show")
		return err
	}

	logger.Info("Successfully fetched show", slog.Any("Show", show))
	schemas.SendSucess(c, "Get Show", show)
	return nil
}
