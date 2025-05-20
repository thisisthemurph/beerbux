package history

import (
	"github.com/google/uuid"
	"time"
)

type SessionHistoryResponse struct {
	SessionID uuid.UUID             `json:"sessionId"`
	Events    []SessionHistoryEvent `json:"events"`
}

type SessionHistoryEvent struct {
	ID        int32       `json:"id"`
	MemberID  uuid.UUID   `json:"memberId"`
	EventType string      `json:"eventType"`
	EventData interface{} `json:"eventData,omitempty"`
	CreatedAt time.Time   `json:"createdAt"`
}

type TransactionCreatedEventData struct {
	TransactionID uuid.UUID         `json:"transactionId"`
	Lines         []TransactionLine `json:"lines"`
}

type MemberAddedEventData struct {
	MemberID uuid.UUID `json:"memberId"`
}

type MemberPromotedToAdminEventData struct {
	MemberID uuid.UUID `json:"memberId"`
}

type MemberDemotedFromAdminEventData struct {
	MemberID uuid.UUID `json:"memberId"`
}

type MemberRemovedEventData struct {
	MemberID uuid.UUID `json:"memberId"`
}

type TransactionLine struct {
	MemberID uuid.UUID `json:"memberId"`
	Amount   float64   `json:"amount"`
}
