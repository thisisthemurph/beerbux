package history

import (
	"beerbux/internal/session/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"log/slog"
)

type SessionHistoryReader interface {
	GetSessionHistory(ctx context.Context, sessionID uuid.UUID) (SessionHistoryResponse, error)
}

type SessionHistoryWriter interface {
	CreateSessionOpenedEvent(ctx context.Context, sessionID, memberID uuid.UUID) error
	CreateSessionClosedEvent(ctx context.Context, sessionID, memberID uuid.UUID) error
	CreateMemberAddedEvent(ctx context.Context, sessionID, memberID, performedByMemberId uuid.UUID) error
	CreateMemberRemovedEvent(ctx context.Context, sessionID, memberID, performedByMemberId uuid.UUID) error
	CreateMemberLeftEvent(ctx context.Context, sessionID, memberID uuid.UUID) error
	CreateMemberPromotedToAdminEvent(ctx context.Context, sessionID, memberID, performedByMemberId uuid.UUID) error
	CreateMemberDemotedFromAdminEvent(ctx context.Context, sessionID, memberID, performedByMemberId uuid.UUID) error
	CreateTransactionCreatedEvent(ctx context.Context, sessionID, performedByMemberId uuid.UUID, transactionLines TransactionHistory) error
}

type SessionHistoryService struct {
	Queries *db.Queries
	logger  *slog.Logger
}

func NewSessionHistoryService(queries *db.Queries, logger *slog.Logger) *SessionHistoryService {
	return &SessionHistoryService{
		Queries: queries,
		logger:  logger,
	}
}

func (r *SessionHistoryService) GetSessionHistory(ctx context.Context, sessionID uuid.UUID) (SessionHistoryResponse, error) {
	events, err := r.Queries.GetSessionHistory(ctx, sessionID)
	if err != nil {
		return SessionHistoryResponse{}, fmt.Errorf("failed to get sessio history events: %w", err)
	}

	response := SessionHistoryResponse{
		SessionID: sessionID,
		Events:    make([]SessionHistoryEvent, 0, len(events)),
	}

	for _, e := range events {
		eventData, err := r.parseEventJSON(e.EventType, e.EventData)
		if err != nil {
			r.logger.Error("failed to parse event data", "session", sessionID, "error", err)
		}

		response.Events = append(response.Events, SessionHistoryEvent{
			ID:        e.ID,
			MemberID:  e.MemberID,
			EventType: e.EventType,
			EventData: eventData,
			CreatedAt: e.CreatedAt,
		})
	}

	return response, nil
}

func (r *SessionHistoryService) parseEventJSON(eventType string, data pqtype.NullRawMessage) (interface{}, error) {
	if !data.Valid {
		return nil, fmt.Errorf("event type %s NillRawMessage is null", eventType)
	}

	switch eventType {
	case EventTransactionCreated:
	case EventMemberAdded:
		var eventData MemberAddedEventData
		if err := json.Unmarshal(data.RawMessage, &eventData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s event data: %w", eventData, err)
		}
		return eventData, nil
	case EventMemberRemoved:
		var eventData MemberRemovedEventData
		if err := json.Unmarshal(data.RawMessage, &eventData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s event data: %w", eventData, err)
		}
		return eventData, nil
	case EventMemberLeft, EventSessionClosed, EventSessionOpened:
		return nil, nil
	case EventMemberPromotedToAdmin:
		return nil, nil
	case EventMemberDemotedFromAdmin:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown event type %s", eventType)
	}
}

func (r *SessionHistoryService) CreateSessionOpenedEvent(ctx context.Context, sessionID, memberID uuid.UUID) error {
	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  memberID,
		EventType: EventSessionOpened,
		EventData: newNilNullRawMessage(),
	})
}

func (r *SessionHistoryService) CreateSessionClosedEvent(ctx context.Context, sessionID, memberID uuid.UUID) error {
	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  memberID,
		EventType: EventSessionClosed,
		EventData: newNilNullRawMessage(),
	})
}

func (r *SessionHistoryService) CreateMemberAddedEvent(ctx context.Context, sessionID, memberID, performedByMemberId uuid.UUID) error {
	eventData := MemberAddedEventData{
		MemberID: memberID,
	}

	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  performedByMemberId,
		EventType: EventMemberAdded,
		EventData: newNullRawMessage(eventData),
	})
}

func (r *SessionHistoryService) CreateMemberRemovedEvent(ctx context.Context, sessionID, memberID, performedByMemberId uuid.UUID) error {
	eventData := MemberRemovedEventData{
		MemberID: memberID,
	}

	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  performedByMemberId,
		EventType: EventMemberRemoved,
		EventData: newNullRawMessage(eventData),
	})
}

func (r *SessionHistoryService) CreateMemberLeftEvent(ctx context.Context, sessionID, memberID uuid.UUID) error {
	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  memberID,
		EventType: EventMemberLeft,
		EventData: newNilNullRawMessage(),
	})
}

type MemberPromotedOrDemotedEventData struct {
	MemberID uuid.UUID `json:"member_id"`
}

func (r *SessionHistoryService) CreateMemberPromotedToAdminEvent(ctx context.Context, sessionID, memberID, performedByMemberId uuid.UUID) error {
	eventData := MemberPromotedOrDemotedEventData{
		MemberID: memberID,
	}

	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  performedByMemberId,
		EventType: EventMemberPromotedToAdmin,
		EventData: newNullRawMessage(eventData),
	})
}

func (r *SessionHistoryService) CreateMemberDemotedFromAdminEvent(ctx context.Context, sessionID, memberID, performedByMemberId uuid.UUID) error {
	eventData := MemberPromotedOrDemotedEventData{
		MemberID: memberID,
	}

	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  performedByMemberId,
		EventType: EventMemberDemotedFromAdmin,
		EventData: newNullRawMessage(eventData),
	})
}

type TransactionHistoryLine struct {
	MemberID uuid.UUID `json:"member_id"`
	Amount   float64   `json:"amount"`
}

type TransactionHistory struct {
	TransactionID uuid.UUID                `json:"transaction_id"`
	Lines         []TransactionHistoryLine `json:"lines"`
}

func (r *SessionHistoryService) CreateTransactionCreatedEvent(
	ctx context.Context,
	sessionID,
	performedByMemberId uuid.UUID,
	transactionLines TransactionHistory,
) error {
	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  performedByMemberId,
		EventType: EventTransactionCreated,
		EventData: newNullRawMessage(transactionLines),
	})

}

func newNullRawMessage(v interface{}) pqtype.NullRawMessage {
	data, _ := json.Marshal(v)
	return pqtype.NullRawMessage{
		RawMessage: data,
		Valid:      true,
	}
}

func newNilNullRawMessage() pqtype.NullRawMessage {
	return pqtype.NullRawMessage{
		RawMessage: nil,
		Valid:      true,
	}
}
