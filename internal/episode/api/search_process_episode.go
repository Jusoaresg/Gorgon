package api

import (
	"log/slog"
	"net/http"

	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Search Process Episode
// @Description Search and Process the episode ( Automatic Search )
// @Tags Database/Show/Episode
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/episode/search [post]
func (h *Handler) SearchProcessEpisode(c echo.Context) error {
	h.Logger.Info("Received request to Search Process Episode", slog.String("endpoint", "/database/show/episode/search"), slog.String("method", c.Request().Method))

	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 500, "Failed to bind body")
		return err
	}

	err := h.EpisodeSearchSvc.ProcessSingleEpisode(int(request.Id))
	if err != nil {
		schemas.SendError(c, http.StatusInternalServerError, "Failed to process episode", nil)
		return err
	}

	schemas.SendSuccess(c, "Process Single Episode", request)
	return nil
}
