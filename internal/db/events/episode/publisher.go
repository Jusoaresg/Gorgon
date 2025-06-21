package episode

import (
	"gorgon/internal/event"
	"gorgon/pkg/handler"
)

func EmitEpisodeTrackingUpdatedEvent(episodeID int64, tracking string) {
	msg := EpisodeTrackingUpdatedPayload{
		Type:      string(event.EpisodeTrackingUpdated),
		EpisodeID: episodeID,
		Tracking:  tracking,
	}
	handler.SendWebSocketMessage(msg)
}
