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
