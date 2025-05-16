package history

import (
	"beerbux/internal/session/db"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"time"
)

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
}

func NewSessionHistoryService(queries *db.Queries) *SessionHistoryService {
	return &SessionHistoryService{
		Queries: queries,
	}
}

type SessionHistoryEvent struct {
	ID        int32
	MemberID  uuid.UUID
	EventType string
	EventData *json.RawMessage
	CreatedAt time.Time
}

type SessionHistoryResponse struct {
	SessionID uuid.UUID
	Events    []SessionHistoryEvent
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

type MemberAddedEventData struct {
	MemberID uuid.UUID `json:"member_id"`
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

type MemberRemovedEventData struct {
	MemberID uuid.UUID `json:"member_id"`
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
