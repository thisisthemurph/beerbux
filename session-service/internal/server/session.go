package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrSessionNotFound = status.Error(codes.NotFound, "session not found")
	ErrUserNotFound    = status.Error(codes.NotFound, "user not found")
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ssn, err := s.sessionRepository.GetSession(ctx, r.SessionId)
	if err != nil {
		s.logger.Error("failed to get session", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, status.Errorf(codes.Internal, "failed to get session: %v", err)
	}

	members, err := s.sessionRepository.ListMembers(ctx, r.SessionId)
	if err != nil {
		s.logger.Error("failed to list members", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list members: %v", err)
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

// ListSessionsForUser returns the sessions for a user containing the associated members.
func (s *SessionServer) ListSessionsForUser(ctx context.Context, r *sessionpb.ListSessionsForUserRequest) (*sessionpb.ListSessionsForUserResponse, error) {
	if err := validateListSessionsForUserRequest(r); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Fetch sessions with members from the database.
	rows, err := s.sessionRepository.ListSessionsForUser(ctx, r.UserId)
	if err != nil {
		s.logger.Error("failed to list sessions", "user_id", r.UserId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list sessions: %v", err)
	}

	sessionsMap := make(map[string]*sessionpb.SessionResponse, len(rows))

	for _, row := range rows {
		ssn, exists := sessionsMap[row.ID]
		if !exists {
			ssn = &sessionpb.SessionResponse{
				SessionId: row.ID,
				Name:      row.Name,
				IsActive:  row.IsActive,
				Members:   make([]*sessionpb.SessionMember, 0, 4),
			}
			sessionsMap[row.ID] = ssn
		}

		if row.MemberID != "" {
			ssn.Members = append(ssn.Members, &sessionpb.SessionMember{
				UserId:   row.MemberID,
				Name:     row.MemberName,
				Username: row.MemberUsername,
			})
		}
	}

	sessions := make([]*sessionpb.SessionResponse, 0, len(sessionsMap))
	for _, ssn := range sessionsMap {
		sessions = append(sessions, ssn)
	}

	return &sessionpb.ListSessionsForUserResponse{
		Sessions: sessions,
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to begin tx: %v", err)
	}
	defer tx.Rollback()

	qtx := s.sessionRepository.WithTx(tx)

	// Add the session to the sessions table
	ssn, err := qtx.CreateSession(ctx, session.CreateSessionParams{
		ID:   uuid.New().String(),
		Name: r.Name,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v", err)
	}

	// Add the user to the members table.
	err = qtx.UpsertMember(ctx, session.UpsertMemberParams{
		ID:       user.UserId,
		Name:     user.Name,
		Username: user.Username,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upsert member: %v", err)
	}

	// Join the member to the sessions in the session_members table.
	err = qtx.AddSessionMember(ctx, session.AddSessionMemberParams{
		SessionID: ssn.ID,
		MemberID:  user.UserId,
		IsOwner:   true,
	})

	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit tx: %v", err)
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
		return nil, status.Errorf(codes.Internal, "failed to begin tx: %v", err)
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
		if errors.Is(sessionErr, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, status.Errorf(codes.Internal, "failed fetching session: %v", sessionErr)
	}

	if memberErr != nil && !errors.Is(memberErr, sql.ErrNoRows) {
		return nil, status.Errorf(codes.Internal, "error fetching member %q from database: %v", r.UserId, memberErr)
	}

	// If the user is not in the members table, fetch the user from the user service.
	if !userInMembersTable {
		u, err := s.userClient.GetUser(ctx, &userpb.GetUserRequest{
			UserId: r.UserId,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrUserNotFound
			}
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
			return nil, status.Errorf(codes.Internal, "failed to upsert member: %v", err)
		}
	}

	// Associate the member with the session.
	err = qtx.AddSessionMember(ctx, session.AddSessionMemberParams{
		SessionID: r.SessionId,
		MemberID:  r.UserId,
		IsOwner:   false,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add member to session: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit tx: %v", err)
	}

	if err := s.sessionMemberAddedPublisher.Publish(r.SessionId, r.UserId); err != nil {
		s.logger.Error("failed to publish session member added event", "error", err)
	}

	return &sessionpb.EmptyResponse{}, nil
}
