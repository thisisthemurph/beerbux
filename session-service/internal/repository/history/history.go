package history

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/thisisthemurph/beerbux/session-service/protos/historypb"
	"google.golang.org/protobuf/types/known/anypb"
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
	GetBySessionID(ctx context.Context, sessionID string) (*historypb.SessionHistoryResponse, error)
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

func (r *SQLiteHistoryRepository) GetBySessionID(ctx context.Context, sessionID string) (*historypb.SessionHistoryResponse, error) {
	events, err := r.queries.GetSessionHistory(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	response := &historypb.SessionHistoryResponse{
		SessionId: sessionID,
		Events:    make([]*historypb.SessionHistoryEvent, 0, len(events)),
	}

	for _, e := range events {
		eventType := NewEventType(e.EventType)
		data, err := parseEvent(eventType, e.EventData)
		if err != nil {
			return nil, err
		}

		response.Events = append(response.Events, &historypb.SessionHistoryEvent{
			Id:        e.ID,
			MemberId:  e.MemberID,
			EventType: e.EventType,
			EventData: data,
			CreatedAt: e.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func parseEvent(eventType EventType, b []byte) (*anypb.Any, error) {
	switch eventType {
	case EventTransactionCreated:
		var transactionCreatedEvent *historypb.TransactionCreatedEventData
		if err := json.Unmarshal(b, &transactionCreatedEvent); err != nil {
			return nil, err
		}
		a, err := anypb.New(transactionCreatedEvent)
		if err != nil {
			return nil, err
		}
		return a, err
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
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
