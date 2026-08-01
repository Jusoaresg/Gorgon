package api

import (
	"log/slog"
	"strconv"

	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Search Process Season Episodes
// @Description Search and process all aired episodes in a season (Season Force Search)
// @Tags Database/Show/Episode
// @Produce json
// @Param id path int true "Season ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/episode/search/season/:id [post]
func (h *Handler) SearchProcessSeasonEpisodes(c echo.Context) error {
	h.Logger.Info("Received request to Search Process Season Episodes",
		slog.String("endpoint", "/database/show/episode/search/season/:id"),
		slog.String("method", c.Request().Method))

	seasonID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		schemas.SendError(c, 400, "Invalid season ID", SearchResponse{
			ToastMessage: "Failed to start season search",
		})
		return err
	}

	h.EpisodeSearchSvc.ProcessSeasonEpisodes(seasonID)

	schemas.SendSuccess(c, "Season search started", SearchResponse{
		ToastMessage: "Season search started",
	})
	return nil
}
