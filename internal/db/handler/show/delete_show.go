package show

import (
	"errors"
	"gorgon/config"
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"
	"log/slog"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Delete Show
// @Description Delete Show from list
// @Tags Database/Show
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show [delete]
func DeleteShow(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Delete Anime", slog.String("endpoint", "/api/v1/database/show"))

	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to bind request body")
		return err
	}

	id := request.Id
	if request.Id == 0 {
		logger.Error("Invalid ID", slog.String("error", "ID cannot be 0"))
		schemas.SendError(c, 400, "ID cannot be 0")
		return errors.New("Id cannot be 0")
	}

	show := model.Show{}

	baseService := services.NewBaseService()
	if err := baseService.Get(&show, id); err != nil {
		logger.Error("Error while retrieving show data from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while retrieving anime data from database")
		return err
	}

	if err := baseService.Delete(show.ID, model.Show{}); err != nil {
		logger.Error("Error while deleting show from databases", slog.Int("id", request.Id), slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while deleting anime from database")
		return err
	}

	schemas.SendSucess(c, "DeleteAnime", show)
	return nil
}
