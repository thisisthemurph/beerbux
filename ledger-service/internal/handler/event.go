package handler

import (
	"github.com/google/uuid"
	"time"
)

type EventMetadata struct {
	Event         string    `json:"event"`
	Version       string    `json:"version"`
	Timestamp     time.Time `json:"timestamp"`
	Source        string    `json:"source,omitempty"`
	CorrelationID uuid.UUID `json:"correlation_id,omitempty"`
}

type TransactionCreatedMemberAmount struct {
	UserID uuid.UUID `json:"user_id"`
	Amount float64   `json:"amount"`
}

type TransactionCreatedEventData struct {
	TransactionID uuid.UUID                        `json:"transaction_id"`
	CreatorID     uuid.UUID                        `json:"creator_id"`
	SessionID     uuid.UUID                        `json:"session_id"`
	MemberAmounts []TransactionCreatedMemberAmount `json:"member_amounts"`
}

type TransactionCreatedEvent struct {
	EventMetadata
	Data TransactionCreatedEventData `json:"data"`
}
