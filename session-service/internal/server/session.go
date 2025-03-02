package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
)

type SessionServer struct {
	sessionpb.UnimplementedSessionServer
	TX
	sessionRepository *session.Queries
	logger            *slog.Logger
	userClient        userpb.UserClient
}

func NewSessionServer(db *sql.DB, sessionRepository *session.Queries, userClient userpb.UserClient, logger *slog.Logger) *SessionServer {
	return &SessionServer{
		TX:                db,
		sessionRepository: sessionRepository,
		userClient:        userClient,
		logger:            logger,
	}
}

func (s *SessionServer) GetSession(ctx context.Context, r *sessionpb.GetSessionRequest) (*sessionpb.SessionResponse, error) {
	if err := validateGetSessionRequest(r); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	ssn, err := s.sessionRepository.GetSession(ctx, r.SessionId)
	if err != nil {
		s.logger.Error("failed to get session", "error", err)
		return nil, fmt.Errorf("failed to get session %v: %w", r.SessionId, err)
	}

	return &sessionpb.SessionResponse{
		SessionId: ssn.ID,
		Name:      ssn.Name,
		IsActive:  ssn.IsActive,
	}, nil
}

// CreateSession creates a new session in the sessions table.
// The creating user is also added as a member in the session_members table.
func (s *SessionServer) CreateSession(ctx context.Context, r *sessionpb.CreateSessionRequest) (*sessionpb.SessionResponse, error) {
	if err := validateCreateSessionRequest(r); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	user, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{
		UserId: r.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := s.sessionRepository.WithTx(tx)

	// Add the session to the sessions table
	ssn, err := qtx.CreateSession(ctx, session.CreateSessionParams{
		ID:   uuid.New().String(),
		Name: r.Name,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Add the user to the members table.
	err = qtx.UpsertMember(ctx, session.UpsertMemberParams{
		ID:       user.UserId,
		Name:     user.Name,
		Username: user.Username,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to upsert member: %w", err)
	}

	// Join the member to the sessions in the session_members table.
	err = qtx.AddSessionMember(ctx, session.AddSessionMemberParams{
		SessionID: ssn.ID,
		MemberID:  user.UserId,
		IsOwner:   true,
	})

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &sessionpb.SessionResponse{
		SessionId: ssn.ID,
		Name:      ssn.Name,
		IsActive:  true,
	}, nil
}
