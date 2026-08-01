package episode

import (
	"github.com/jusoaresg/gorgon/internal/event"
	"github.com/jusoaresg/gorgon/pkg/handler"
)

func EmitEpisodeTrackingUpdatedEvent(episodeID int64, tracking string) {
	msg := EpisodeTrackingUpdatedPayload{
		Type:      string(event.EpisodeTrackingUpdated),
		EpisodeID: episodeID,
		Tracking:  tracking,
	}
	handler.SendWebSocketMessage(msg)
}

func EmitEpisodeSearchFinishedEvent(episodeID int64, season, number int, name, showName, result, message string) {
	msg := EpisodeSearchFinishedPayload{
		Type:      string(event.EpisodeSearchFinished),
		EpisodeID: episodeID,
		Season:    season,
		Number:    number,
		Name:      name,
		ShowName:  showName,
		Result:    result,
		Message:   message,
	}
	handler.SendWebSocketMessage(msg)
}
