package show

import (
	"fmt"
	"gorgon/assets/templates/components/add_show"
	"gorgon/config"
	"gorgon/external/tvmaze/service"
	showManager "gorgon/internal/db/service"
	"gorgon/internal/templates"
	"gorgon/pkg/schemas"
	"gorgon/pkg/services"
	"log/slog"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Add Show
// @Description Add Show to List
// @Tags Database/Show
// @Accept json
// @Produce json
// @Param request body schemas.IdRequest true "Request Body"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show [post]
func AddShowToList(c echo.Context) error {
	logger := config.GetLogger()
	logger.Info("Received request to Add Show To List", slog.String("endpoint", "/api/v1/database/show"), slog.String("method", "POST"))

	var request schemas.IdRequest

	isHTMX := c.Request().Header.Get("HX-Request") == "true"
	if isHTMX {
		id, err := strconv.ParseInt(c.FormValue("id"), 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		request.Id = int(id)
	} else {
		if err := c.Bind(&request); err != nil {
			logger.Error("Failed to bind body request")
			schemas.SendError(c, 500, "Failed to bind body request")
			return err
		}
	}

	tvMazeService := service.NewTvMazeSearchService(logger)
	showManagerService := showManager.NewShowManagerService(logger)

	showDto, err := tvMazeService.SearchByTvMazeId(request.Id)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	episodesDto, err := showManagerService.GetEpisodes(showDto.ShowID)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	seasonsDto, err := showManagerService.GetSeasons(showDto.ShowID)
	if err != nil {
		schemas.SendError(c, 500, err.Error())
		return err
	}

	show := showDto.ToModel(episodesDto, seasonsDto)

	baseService := services.NewBaseService()
	if err := baseService.Add(&show); err != nil {
		logger.Error("Failed to add anime to database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Failed to add anime to database")
		return err
	}

	if isHTMX {
		// return c.Render(200, "", "")
		return templates.Render(c, add_show.ShowCard(*showDto, true))
		// @add_show.ShowCard(show.Show, addedShows[show.Show.ShowID])
	}

	schemas.SendSucess(c, "Add Show To List", &show)
	return nil
}
