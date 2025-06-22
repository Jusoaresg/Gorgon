package show

import (
	"errors"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/show/model"
	"github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"log/slog"

	"github.com/jmoiron/sqlx"
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
	return deleteShowHandler(c, config.GetSQLite())
}

func deleteShowHandler(c echo.Context, db *sqlx.DB) error {
	logger := config.GetLogger()
	logger.Info("Received request to Delete Show", slog.String("endpoint", "/api/v1/database/show"))

	var request schemas.IdRequest

	if err := c.Bind(&request); err != nil {
		logger.Error("Failed to bind request body", slog.String("error", err.Error()))
		schemas.SendError(c, 400, "Failed to bind request body")
		return err
	}

	if request.Id <= 0 {
		logger.Error("Invalid ID", slog.String("error", "ID must be greater than 0"))
		schemas.SendError(c, 400, "ID must be greater than 0")
		return errors.New("ID must be greater than 0")
	}

	show := model.Show{}

	showRepo := repository.NewShowRepository(db)
	err := showRepo.DeleteById(request.Id)
	if err != nil {
		if errors.Is(err, repository.ErrShowNotFound) {
			schemas.SendError(c, 404, "Show not found")
			return err
		}
		schemas.SendError(c, 500, "Internal server error")
		return err
	}

	schemas.SendSuccess(c, "DeleteShow", show)
	return nil
}
