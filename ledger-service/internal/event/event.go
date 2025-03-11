package event

import (
	"github.com/google/uuid"
	"time"
)

type Metadata struct {
	Event         string    `json:"event"`
	Version       string    `json:"version"`
	Timestamp     time.Time `json:"timestamp"`
	Source        string    `json:"source,omitempty"`
	CorrelationID string    `json:"correlation_id,omitempty"`
}

func NewMetadata(event, version, correlationID string) Metadata {
	return Metadata{
		Event:         event,
		Version:       version,
		Timestamp:     time.Now(),
		Source:        "ledger-service",
		CorrelationID: correlationID,
	}
}

type TransactionCreatedMemberAmount struct {
	UserID uuid.UUID `json:"user_id"`
	Amount float64   `json:"amount"`
}

type TransactionCreatedEvent struct {
	TransactionID uuid.UUID                        `json:"transaction_id"`
	CreatorID     uuid.UUID                        `json:"creator_id"`
	SessionID     uuid.UUID                        `json:"session_id"`
	MemberAmounts []TransactionCreatedMemberAmount `json:"member_amounts"`
}
