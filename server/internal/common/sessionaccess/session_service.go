package sessionaccess

import (
	"beerbux/internal/common/sessionaccess/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/thisisthemurph/fn"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrMemberNotFound  = errors.New("member not found")
)

type SessionReader interface {
	// GetSessionByID returns the session for the given ID, including the members and transaction lines.
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*SessionWithTransactions, error)
	// GetSessionDetails returns the basic data for the session, including the members of the session.
	GetSessionDetails(ctx context.Context, sessionID uuid.UUID) (*Session, error)
	// GetSessionMember returns the SessionMember for the given session and member ID or an error if the member or session does not exist.
	GetSessionMember(ctx context.Context, sessionID, memberID uuid.UUID) (*SessionMember, error)
	// UserIsMemberOfSession returns a bool indicating if the session includes a member with the given ID.
	UserIsMemberOfSession(ctx context.Context, sessionID, memberID uuid.UUID) (bool, error)
}

type SessionService struct {
	queries *db.Queries
}

func NewSessionService(queries *db.Queries) SessionReader {
	return &SessionService{
		queries: queries,
	}
}

func (s *SessionService) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*SessionWithTransactions, error) {
	session, err := s.queries.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session %s: %w", sessionID, err)
	}

	lines, err := s.queries.GetSessionTransactionLines(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction lines for session %s: %w", sessionID, err)
	}

	members, err := s.queries.ListSessionMembers(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list session %s members: %w", sessionID, err)
	}

	result := &SessionWithTransactions{
		Session: Session{
			ID:       session.ID,
			Name:     session.Name,
			IsActive: session.IsActive,
			Members:  s.buildSessionMembers(members),
		},
		Transactions: s.buildSessionTransactions(lines),
		Total:        session.Total,
	}

	return result, nil
}

func (s *SessionService) GetSessionDetails(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	session, err := s.queries.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session %s: %w", sessionID, err)
	}

	members, err := s.queries.ListSessionMembers(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list session %s members: %w", sessionID, err)
	}

	result := &Session{
		ID:       session.ID,
		Name:     session.Name,
		IsActive: session.IsActive,
		Members:  s.buildSessionMembers(members),
	}

	return result, nil
}

func (s *SessionService) GetSessionMember(ctx context.Context, sessionID, memberID uuid.UUID) (*SessionMember, error) {
	m, err := s.queries.GetSessionMember(ctx, db.GetSessionMemberParams{
		SessionID: sessionID,
		MemberID:  memberID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, fmt.Errorf("error fetching session member: %w", err)
	}

	return &SessionMember{
		ID:        m.ID,
		Name:      m.Name,
		Username:  m.Username,
		IsAdmin:   m.IsAdmin,
		IsDeleted: m.IsDeleted,
	}, nil
}

func (s *SessionService) UserIsMemberOfSession(ctx context.Context, sessionID, memberID uuid.UUID) (bool, error) {
	_, err := s.queries.GetSessionMember(ctx, db.GetSessionMemberParams{
		SessionID: sessionID,
		MemberID:  memberID,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to deyermine if user is member of session: %w", err)
	}

	return true, nil
}

func (s *SessionService) buildSessionMembers(members []db.ListSessionMembersRow) []SessionMember {
	return fn.Map(members, func(m db.ListSessionMembersRow) SessionMember {
		return SessionMember{
			ID:        m.ID,
			Name:      m.Name,
			Username:  m.Username,
			IsAdmin:   m.IsAdmin,
			IsDeleted: m.IsDeleted,
		}
	})
}

func (s *SessionService) buildSessionTransactions(lines []db.GetSessionTransactionLinesRow) []SessionTransaction {
	if len(lines) == 0 {
		return make([]SessionTransaction, 0)
	}

	transactionMap := make(map[uuid.UUID]*SessionTransaction)
	for _, line := range lines {
		if _, exists := transactionMap[line.TransactionID]; !exists {
			transactionMap[line.TransactionID] = &SessionTransaction{
				ID:        line.TransactionID,
				UserID:    line.CreatorID,
				Total:     0,
				Lines:     make([]SessionTransactionLine, 0),
				CreatedAt: line.CreatedAt,
			}
		}

		transactionMap[line.TransactionID].Total += line.Amount
		transactionMap[line.TransactionID].Lines = append(transactionMap[line.TransactionID].Lines, SessionTransactionLine{
			UserID: line.MemberID,
			Amount: line.Amount,
		})
	}

	transactions := make([]SessionTransaction, 0, len(transactionMap))
	for _, transaction := range transactionMap {
		transactions = append(transactions, *transaction)
	}

	return transactions
}
