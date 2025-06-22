package episode

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/db/events/episode"
	"github.com/jusoaresg/gorgon/internal/db/model"
	"github.com/jusoaresg/gorgon/internal/db/repository"
	"github.com/jusoaresg/gorgon/internal/paths"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/jusoaresg/gorgon/utils"
	"log/slog"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// @BasePath /api/v1

// @Summary Delete Downloaded Episode
// @Description Delete Downloaded Episode
// @Tags Database/Episodes
// @Produce json
// @Param id path int true "Episode ID"
// @Success 200 {object} schemas.DefaultResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /database/show/episode/{id} [delete]
func DeleteDownloadedEpisode(c echo.Context) error {
	return deleteDownloadedEpisodeHandler(c, config.GetSQLite())
}

func deleteDownloadedEpisodeHandler(c echo.Context, db *sqlx.DB) error {
	logger := config.GetLogger()

	logger.Info("Received request to Delete Downloaded Episode", slog.String("endpoint", "/database/show/episode/:id"), slog.String("method", "delete"))

	cfg, err := config.LoadConfig()
	if err != nil {
		schemas.SendError(c, 400, "Error while loading config file")
		return err
	}

	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		schemas.SendError(c, 400, "Error while converting id to int")
		return err
	}

	epRepo := repository.NewEpisodeRepository(db)
	epContentRepo := repository.NewEpisodeContentRepository()
	showRepo := repository.NewShowRepository(db)

	ep, err := epRepo.GetByID(idInt64)
	if err != nil {
		logger.Error("Error while fetching episode from database to delete", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching episode to delete")
		return err
	}

	episodeContent, err := epContentRepo.GetByEpisodeId(idInt64)
	if err != nil {
		logger.Error("Error while fetching episode content from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching episode content")
		return err
	}

	show, err := showRepo.GetById(ep.ShowID)
	if err != nil {
		logger.Error("Error while fetching show from database", slog.String("error", err.Error()))
		schemas.SendError(c, 500, "Error while fetching show")
		return err
	}

	downloadFolder, err := paths.GetTorrentDownloadFolder()
	if err != nil {
		return err
	}

	filePath, err := utils.DeleteFile(downloadFolder, episodeContent.Name)
	if err != nil {
		logger.Error("Failed to delete downloaded episode", slog.String("error", err.Error()))
		return err
	}
	if err := utils.DeleteSymlink(cfg.ShowsFolder, show.Name, ep, episodeContent); err != nil {
		logger.Error("Failed to delete episode symlink from shows folder", slog.String("error", err.Error()))
		return err
	}

	if err := epContentRepo.DeleteById(episodeContent.ID); err != nil {
		logger.Error("Failed to delete episode content from the database (The File was deleted)", slog.String("error", err.Error()))
	}

	ep.SetNotInstalled()
	err = epRepo.Update(ep)
	if err != nil {
		logger.Error("Failed to set episode to not intalled", slog.String("error", err.Error()))
	}

	episode.EmitEpisodeTrackingUpdatedEvent(ep.ID, model.TrackingSkipped)

	logger.Info("Successfully deleted downloaded episode", slog.String("FilePath", filePath))
	schemas.SendSucess(c, "Delete Downloaded Episode", episodeContent)
	return nil
}
