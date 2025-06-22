package services

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/internal/db/model"
	"github.com/jusoaresg/gorgon/pkg/schemas/dtos"
	"log/slog"
	"time"
)

func ApplyTrackingToEpisodes(episodes *[]dtos.EpisodeDto, tracking string) {
	now := time.Now()
	logger := config.GetLogger()

	for i, ep := range *episodes {
		airStampTime, err := time.Parse(time.RFC3339, ep.AirStamp)
		if err != nil {
			logger.Info("Failed to set episode airStampTime, setting episode tracking to wanted", slog.String("Episode Name", ep.Name), slog.Int("Episode Number", ep.Number))
			(*episodes)[i].Tracking = model.TrackingWanted
			continue
		}
		switch tracking {
		case "all":
			(*episodes)[i].Tracking = model.TrackingWanted
		case "future":
			if airStampTime.After(now) {
				(*episodes)[i].Tracking = model.TrackingWanted
			} else {
				(*episodes)[i].Tracking = model.TrackingSkipped
			}
		case "none":
			(*episodes)[i].Tracking = model.TrackingSkipped
		}
	}
}
