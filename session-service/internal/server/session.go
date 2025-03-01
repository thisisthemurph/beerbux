package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"log/slog"
)

type SessionServer struct {
	sessionpb.UnimplementedSessionServer
	sessionRepository *session.Queries
	logger            *slog.Logger
}

func NewSessionServer(sessionRepository *session.Queries, logger *slog.Logger) *SessionServer {
	return &SessionServer{
		sessionRepository: sessionRepository,
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

func (s *SessionServer) CreateSession(ctx context.Context, r *sessionpb.CreateSessionRequest) (*sessionpb.SessionResponse, error) {
	if err := validateCreateSessionRequest(r); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	ssn, err := s.sessionRepository.CreateSession(ctx, session.CreateSessionParams{
		ID:      uuid.New().String(),
		Name:    r.Name,
		OwnerID: r.UserId,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &sessionpb.SessionResponse{
		SessionId: ssn.ID,
		Name:      ssn.Name,
		IsActive:  true,
	}, nil
}
