package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrSessionNotFound                   = status.Error(codes.NotFound, "session not found")
	ErrUserNotFound                      = status.Error(codes.NotFound, "user not found")
	ErrSessionMustHaveAtLeastOneAdmin    = status.Error(codes.FailedPrecondition, "session must have at least one admin")
	ErrorSessionMustHaveAtLeastOneMember = status.Error(codes.FailedPrecondition, "session must have at least one member")
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

func (s *SessionServer) GetSession(ctx context.Context, r *sessionpb.GetSessionRequest) (*sessionpb.GetSessionResponse, error) {
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

	lines, err := s.sessionRepository.GetSessionTransactionLines(ctx, r.SessionId)
	if err != nil {
		s.logger.Error("failed to get session transaction lines", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get session transaction lines: %v", err)
	}

	members, err := s.sessionRepository.ListSessionMembers(ctx, r.SessionId)
	if err != nil {
		s.logger.Error("failed to list members", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list members: %v", err)
	}

	sessionMembers := make([]*sessionpb.SessionMember, 0, len(members))
	for _, m := range members {
		sessionMembers = append(sessionMembers, &sessionpb.SessionMember{
			UserId:    m.ID,
			Name:      m.Name,
			Username:  m.Username,
			IsOwner:   m.IsOwner,
			IsAdmin:   m.IsAdmin,
			IsDeleted: m.IsDeleted,
		})
	}

	transactionMap := make(map[string]*sessionpb.SessionTransaction)
	for _, line := range lines {
		if _, exists := transactionMap[line.TransactionID]; !exists {
			transactionMap[line.TransactionID] = &sessionpb.SessionTransaction{
				TransactionId: line.TransactionID,
				UserId:        line.CreatorID,
				Total:         0,
				Lines:         make([]*sessionpb.SessionTransactionLine, 0),
				CreatedAt:     line.CreatedAt.Format(time.RFC3339),
			}
		}

		transactionMap[line.TransactionID].Total += line.Amount
		transactionMap[line.TransactionID].Lines = append(transactionMap[line.TransactionID].Lines, &sessionpb.SessionTransactionLine{
			UserId: line.MemberID,
			Amount: line.Amount,
		})
	}

	result := &sessionpb.GetSessionResponse{
		SessionId:    ssn.ID,
		Name:         ssn.Name,
		IsActive:     ssn.IsActive,
		Members:      sessionMembers,
		Total:        ssn.Total,
		Transactions: make([]*sessionpb.SessionTransaction, 0, len(transactionMap)),
	}

	for _, transaction := range transactionMap {
		result.Transactions = append(result.Transactions, transaction)
	}

	return result, nil
}

// ListSessionsForUser returns the sessions for a user containing the associated members.
func (s *SessionServer) ListSessionsForUser(ctx context.Context, r *sessionpb.ListSessionsForUserRequest) (*sessionpb.ListSessionsForUserResponse, error) {
	if err := validateListSessionsForUserRequest(r); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	rows, err := s.sessionRepository.ListSessionsForUser(ctx, session.ListSessionsForUserParams{
		MemberID: r.UserId,
		PageSize: r.PageSize,
	})
	if err != nil {
		s.logger.Error("failed to list sessions", "user_id", r.UserId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list sessions: %v", err)
	}

	// Maintain order using a slice
	var sessions []*sessionpb.SessionResponse
	sessionsMap := make(map[string]*sessionpb.SessionResponse, len(rows))

	for _, row := range rows {
		ssn, exists := sessionsMap[row.ID]
		if !exists {
			ssn = &sessionpb.SessionResponse{
				SessionId: row.ID,
				Name:      row.Name,
				IsActive:  row.IsActive,
				Members:   make([]*sessionpb.SessionMember, 0, 4),
				Total:     row.TotalAmount,
			}
			sessionsMap[row.ID] = ssn
			sessions = append(sessions, ssn)
		}

		ssn.Members = append(ssn.Members, &sessionpb.SessionMember{
			UserId:   row.MemberID,
			Name:     row.MemberName,
			Username: row.MemberUsername,
		})
	}

	pageToken := ""
	if r.PageSize > 0 && len(sessions) == int(r.PageSize) {
		pageToken = sessions[len(sessions)-1].SessionId
	}

	return &sessionpb.ListSessionsForUserResponse{
		Sessions:      sessions,
		NextPageToken: pageToken,
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
		IsAdmin:   true,
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

func (s *SessionServer) RemoveMemberFromSession(ctx context.Context, r *sessionpb.RemoveMemberFromSessionRequest) (*sessionpb.EmptyResponse, error) {
	memberCount, err := s.sessionRepository.CountSessionMembers(ctx, r.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count session members: %v", err)
	}
	if memberCount <= 1 {
		return nil, ErrorSessionMustHaveAtLeastOneMember
	}

	memberToRemove, err := s.sessionRepository.GetSessionMember(ctx, session.GetSessionMemberParams{
		ID:        r.UserId,
		SessionID: r.SessionId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch member from session: %v", err)
	}

	if memberToRemove.IsAdmin {
		// Ensure the member is not the only admin.
		count, err := s.sessionRepository.CountSessionAdmins(ctx, r.SessionId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to fetch members: %v", err)
		}
		if count <= 1 {
			return nil, ErrSessionMustHaveAtLeastOneAdmin
		}
	}

	err = s.sessionRepository.DeleteSessionMember(ctx, session.DeleteSessionMemberParams{
		SessionID: r.SessionId,
		MemberID:  r.UserId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove member from session: %v", err)
	}
	return nil, nil
}

func (s *SessionServer) UpdateSessionMemberAdminState(ctx context.Context, r *sessionpb.UpdateSessionMemberAdminStateRequest) (*sessionpb.EmptyResponse, error) {
	count, err := s.sessionRepository.CountSessionAdmins(ctx, r.SessionId)
	if err != nil {
		s.logger.Error("failed to count session admins", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to count session admins: %v", err)
	}

	// Prevent the last remaining admin being removed as admin.
	if !r.IsAdmin && count == 1 {
		return nil, ErrSessionMustHaveAtLeastOneAdmin
	}

	err = s.sessionRepository.UpdateSessionMemberAdmin(ctx, session.UpdateSessionMemberAdminParams{
		IsAdmin:   r.IsAdmin,
		SessionID: r.SessionId,
		MemberID:  r.UserId,
	})

	if err != nil {
		s.logger.Error("failed to update session member admin state", "NewState", r.IsAdmin, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update session member admin state: %v", err)
	}

	return nil, nil
}
