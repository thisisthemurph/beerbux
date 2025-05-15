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
		EventType: EventSessionOpened.String(),
		EventData: pqtype.NullRawMessage{
			RawMessage: nil,
			Valid:      true,
		},
	})
}

func (r *SessionHistoryService) CreateSessionClosedEvent(ctx context.Context, sessionID, memberID uuid.UUID) error {
	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  memberID,
		EventType: EventSessionClosed.String(),
		EventData: NewNilNullRawMessage(),
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
		EventType: EventMemberAdded.String(),
		EventData: NewNullRawMessage(eventData),
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
		EventType: EventMemberRemoved.String(),
		EventData: NewNullRawMessage(eventData),
	})
}

func (r *SessionHistoryService) CreateMemberLeftEvent(ctx context.Context, sessionID, memberID uuid.UUID) error {
	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  memberID,
		EventType: EventMemberLeft.String(),
		EventData: NewNilNullRawMessage(),
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
		EventType: EventMemberPromotedToAdmin.String(),
		EventData: NewNullRawMessage(eventData),
	})
}

func (r *SessionHistoryService) CreateMemberDemotedFromAdminEvent(ctx context.Context, sessionID, memberID, performedByMemberId uuid.UUID) error {
	eventData := MemberPromotedOrDemotedEventData{
		MemberID: memberID,
	}

	return r.Queries.CreateSessionHistory(ctx, db.CreateSessionHistoryParams{
		SessionID: sessionID,
		MemberID:  performedByMemberId,
		EventType: EventMemberDemotedFromAdmin.String(),
		EventData: NewNullRawMessage(eventData),
	})
}

func NewNullRawMessage(v interface{}) pqtype.NullRawMessage {
	data, _ := json.Marshal(v)
	return pqtype.NullRawMessage{
		RawMessage: data,
		Valid:      true,
	}
}

func NewNilNullRawMessage() pqtype.NullRawMessage {
	return pqtype.NullRawMessage{
		RawMessage: nil,
		Valid:      true,
	}
}
