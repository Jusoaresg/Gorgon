package api

import (
	"log/slog"

	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Search Process Show Wanted Episodes
// @Description Search and process show wanted episodes ( Automatic Show Search )
// @Tags Database/Show/Episode
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/episode/search/all [post]
func (h *Handler) SearchProcessShowWantedEpisodes(c echo.Context) error {
	h.Logger.Info("Received request to Search Process Show Wanted Episodes", slog.String("endpoint", "/database/show/episode/search/all"), slog.String("method", c.Request().Method))

	var request schemas.IdRequest
	if err := c.Bind(&request); err != nil {
		schemas.SendError(c, 500, "Failed to bind body")
		return err
	}

	h.EpisodeSearchSvc.ProcessShowWantedEpisodes(int(request.Id))

	schemas.SendSuccess(c, "Process Show Wanted Episodes", request)
	return nil
}
