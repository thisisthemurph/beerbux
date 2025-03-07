package publisher

import "time"

type EventMetadata struct {
	Event         string    `json:"event"`
	Version       string    `json:"version"`
	Timestamp     time.Time `json:"timestamp"`
	Source        string    `json:"source,omitempty"`
	CorrelationID string    `json:"correlation_id,omitempty"`
}

func NewEventMetadata(event, version, correlationID string) EventMetadata {
	return EventMetadata{
		Event:         event,
		Version:       version,
		Timestamp:     time.Now(),
		Source:        "transaction-service",
		CorrelationID: correlationID,
	}
}
