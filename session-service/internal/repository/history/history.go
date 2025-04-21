package history

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/thisisthemurph/fn"
	"time"
)

type EventType string

const (
	EventUnknown            EventType = "unknown_event"
	EventTransactionCreated EventType = "transaction_created"
)

func NewEventType(t string) EventType {
	switch t {
	case EventTransactionCreated.String():
		return EventTransactionCreated
	default:
		return EventUnknown
	}
}

func (et EventType) String() string {
	return string(et)
}

type HistoryRepository interface {
	GetBySessionID(ctx context.Context, sessionID string) ([]SessionHistoryResult, error)
	CreateTransactionCreatedEvent(ctx context.Context, sessionID, memberID string, event TransactionCreatedEvent) error
}

type SQLiteHistoryRepository struct {
	queries *Queries
}

func NewHistoryRepository(db *sql.DB) HistoryRepository {
	return &SQLiteHistoryRepository{
		queries: New(db),
	}
}

type SessionHistoryResult struct {
	ID        int64     `json:"id"`
	SessionID string    `json:"sessionId"`
	MemberID  string    `json:"memberId"`
	EventType EventType `json:"eventType"`
	EventData []byte    `json:"eventData"`
	CreatedAt time.Time `json:"createdAt"`
}

func (r *SQLiteHistoryRepository) GetBySessionID(ctx context.Context, sessionID string) ([]SessionHistoryResult, error) {
	events, err := r.queries.GetSessionHistory(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	results := fn.Map(events, func(s SessionHistory) SessionHistoryResult {
		return SessionHistoryResult{
			ID:        s.ID,
			SessionID: s.SessionID,
			MemberID:  s.MemberID,
			EventType: NewEventType(s.EventType),
			EventData: s.EventData,
			CreatedAt: s.CreatedAt,
		}
	})

	return results, nil
}

type TransactionCreatedEventTransactionLine struct {
	MemberID string  `json:"memberId"`
	Amount   float64 `json:"amount"`
}

type TransactionCreatedEvent struct {
	TransactionID string                                   `json:"transactionId"`
	Lines         []TransactionCreatedEventTransactionLine `json:"transactionLines"`
}

func (r *SQLiteHistoryRepository) CreateTransactionCreatedEvent(ctx context.Context, sessionID, memberID string, event TransactionCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction created event: %w", err)
	}

	return r.queries.CreateSessionHistory(ctx, CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  memberID,
		EventType: EventTransactionCreated.String(),
		EventData: data,
	})
}
