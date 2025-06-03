package services

import (
	"gorgon/config"
	"gorgon/internal/db/model"
	"gorgon/pkg/schemas/dtos"
	"log/slog"
	"time"
)

func ApplyTrackingToEpisodes(episodes *[]dtos.EpisodeDto, tracking string) {
	now := time.Now()
	logger := config.GetLogger()

	for i, ep := range *episodes {
		airStampTime, err := time.Parse(time.RFC3339, ep.AirStamp)
		if err != nil {
			logger.Info("Failed to set episode airStampTime, setting episode tracking to wanted", slog.Int("ShowID", ep.ShowId), slog.String("Episode Name", ep.Name), slog.Int("Episode Number", ep.Number))
			(*episodes)[i].Tracking = model.Tracking.Wanted()
			continue
		}
		switch tracking {
		case "all":
			(*episodes)[i].Tracking = model.Tracking.Wanted()
		case "future":
			if airStampTime.After(now) {
				(*episodes)[i].Tracking = model.Tracking.Wanted()
			} else {
				(*episodes)[i].Tracking = model.Tracking.Skipped()
			}
		case "none":
			(*episodes)[i].Tracking = model.Tracking.Skipped()
		}
	}
}
