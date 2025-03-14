package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/thisisthemurph/beerbux/session-service/internal/publisher"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
)

type SessionServer struct {
	sessionpb.UnimplementedSessionServer
	TX

	sessionRepository           *session.Queries
	logger                      *slog.Logger
	userClient                  userpb.UserClient
	sessionMemberAddedPublisher publisher.SessionMemberAddedPublisher
}

func NewSessionServer(
	db *sql.DB,
	sessionRepository *session.Queries,
	userClient userpb.UserClient,
	sessionMemberAddedPublisher publisher.SessionMemberAddedPublisher,
	logger *slog.Logger,
) *SessionServer {
	return &SessionServer{
		TX:                          db,
		sessionRepository:           sessionRepository,
		userClient:                  userClient,
		sessionMemberAddedPublisher: sessionMemberAddedPublisher,
		logger:                      logger,
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

	members, err := s.sessionRepository.ListMembers(ctx, r.SessionId)
	if err != nil {
		s.logger.Error("failed to list members", "error", err)
		return nil, fmt.Errorf("failed to list members for session %v: %w", r.SessionId, err)
	}

	sessionMembers := make([]*sessionpb.SessionMember, 0, len(members))
	for _, m := range members {
		sessionMembers = append(sessionMembers, &sessionpb.SessionMember{
			UserId:   m.ID,
			Name:     m.Name,
			Username: m.Username,
		})
	}

	return &sessionpb.SessionResponse{
		SessionId: ssn.ID,
		Name:      ssn.Name,
		IsActive:  ssn.IsActive,
		Members:   sessionMembers,
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

// AddMemberToSession adds a user to a session.
func (s *SessionServer) AddMemberToSession(ctx context.Context, r *sessionpb.AddMemberToSessionRequest) (*sessionpb.EmptyResponse, error) {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	qtx := s.sessionRepository.WithTx(tx)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	var sessionErr, memberErr error
	var userInMembersTable bool

	// Ensure the session exists.
	go func() {
		defer wg.Done()
		_, localSessionErr := qtx.GetSession(ctx, r.SessionId)
		if localSessionErr != nil {
			s.logger.Error("error fetching session from database", "ID", r.SessionId, "error", sessionErr)
			cancel()
		}

		sessionErr = localSessionErr
	}()

	// Check if the user is already in the members table.
	go func() {
		defer wg.Done()
		_, localMemberErr := qtx.GetMember(ctx, r.UserId)
		if localMemberErr != nil && !errors.Is(localMemberErr, sql.ErrNoRows) {
			s.logger.Error("error fetching member from database", "ID", r.UserId, "error", memberErr)
			cancel()
		}

		memberErr = localMemberErr
		userInMembersTable = localMemberErr == nil
	}()

	wg.Wait()

	if sessionErr != nil {
		return nil, fmt.Errorf("failed fetching session %q from database: %w", r.SessionId, sessionErr)
	}

	if memberErr != nil && !errors.Is(memberErr, sql.ErrNoRows) {
		return nil, fmt.Errorf("error fetching member %q from database: %w", r.UserId, memberErr)
	}

	// If the user is not in the members table, fetch the user from the user service.
	if !userInMembersTable {
		u, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{
			UserId: r.UserId,
		})
		if err != nil {
			s.logger.Error("error fetching user from user service", "ID", r.UserId, "error", err)
			return nil, fmt.Errorf("failed to fetch user: %w", err)
		}

		err = qtx.UpsertMember(ctx, session.UpsertMemberParams{
			ID:       u.UserId,
			Name:     u.Name,
			Username: u.Username,
		})
		if err != nil {
			s.logger.Error("error upserting member", "ID", u.UserId, "error", err)
			return nil, err
		}
	}

	// Associate the member with the session.
	err = qtx.AddSessionMember(ctx, session.AddSessionMemberParams{
		SessionID: r.SessionId,
		MemberID:  r.UserId,
		IsOwner:   false,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if err := s.sessionMemberAddedPublisher.Publish(r.SessionId, r.UserId); err != nil {
		s.logger.Error("failed to publish session member added event", "error", err)
	}

	return &sessionpb.EmptyResponse{}, nil
}
