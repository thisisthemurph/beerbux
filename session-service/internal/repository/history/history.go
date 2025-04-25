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
	EventMemberRemoved      EventType = "member_removed"
	EventMemberLeft         EventType = "member_left"
)

func NewEventType(t string) EventType {
	switch t {
	case EventTransactionCreated.String():
		return EventTransactionCreated
	case EventMemberRemoved.String():
		return EventMemberRemoved
	case EventMemberLeft.String():
		return EventMemberLeft
	default:
		return EventUnknown
	}
}

func (et EventType) String() string {
	return string(et)
}

type HistoryRepository interface {
	GetBySessionID(ctx context.Context, sessionID string) (*historypb.SessionHistoryResponse, error)
	CreateTransactionCreatedEvent(ctx context.Context, sessionID, memberID string, event *historypb.TransactionCreatedEventData) error
	CreateMemberRemovedEvent(ctx context.Context, sessionID, memberID, performedByMemberId string) error
	CreateMemberLeftEvent(ctx context.Context, sessionID, memberID string) error
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
	case EventMemberRemoved:
		var memberRemovedEvent *historypb.MemberRemovedEventData
		if err := json.Unmarshal(b, &memberRemovedEvent); err != nil {
			return nil, err
		}
		a, err := anypb.New(memberRemovedEvent)
		if err != nil {
			return nil, err
		}
		return a, err
	case EventMemberLeft:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}

func (r *SQLiteHistoryRepository) CreateTransactionCreatedEvent(ctx context.Context, sessionID, memberID string, event *historypb.TransactionCreatedEventData) error {
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

func (r *SQLiteHistoryRepository) CreateMemberRemovedEvent(ctx context.Context, sessionID, memberID, performedByMemberId string) error {
	eventData := &historypb.MemberRemovedEventData{
		MemberId: memberID,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal member removed event: %w", err)
	}

	return r.queries.CreateSessionHistory(ctx, CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  performedByMemberId,
		EventType: EventMemberRemoved.String(),
		EventData: data,
	})
}

func (r *SQLiteHistoryRepository) CreateMemberLeftEvent(ctx context.Context, sessionID, memberID string) error {
	return r.queries.CreateSessionHistory(ctx, CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  memberID,
		EventType: EventMemberLeft.String(),
		EventData: nil,
	})
}
