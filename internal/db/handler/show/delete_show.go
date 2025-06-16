package show

import (
	"errors"
	"gorgon/config"
	"gorgon/internal/db/model"
	"gorgon/internal/db/repository"
	"gorgon/pkg/schemas"
	"log/slog"
	"strconv"

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
	logger.Info("Received request to Delete Show", slog.String("endpoint", "/api/v1/database/show"))

	var request schemas.IdRequest

	if err := c.Bind(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to bind request body")
		return err
	}

	idString := strconv.Itoa(request.Id)
	id64, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		logger.Error("Failed to parse id to int64", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "ID cannot be 0")
		return errors.New("Id cannot be 0")
	}

	if request.Id == 0 {
		logger.Error("Invalid ID", slog.String("error", "ID cannot be 0"))
		schemas.SendError(c, 400, "ID cannot be 0")
		return errors.New("Id cannot be 0")
	}

	show := model.Show{}

	showRepo := repository.NewShowRepository()
	showRepo.DeleteById(id64)

	schemas.SendSucess(c, "DeleteShow", show)
	return nil
}
