package event

type EventType string

const (
	EpisodeTrackingUpdated EventType = "EpisodeTrackingUpdated"
	EpisodeSearchFinished  EventType = "EpisodeSearchFinished"
)
