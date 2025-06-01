package sessionaccess

import (
	"beerbux/internal/common/sessionaccess/db"
	sessionErr "beerbux/internal/session/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/thisisthemurph/fn"
)

type SessionReader interface {
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*SessionResponse, error)
}

type SessionService struct {
	queries *db.Queries
}

func NewSessionService(queries *db.Queries) SessionReader {
	return &SessionService{
		queries: queries,
	}
}

func (s *SessionService) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*SessionResponse, error) {
	session, err := s.queries.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sessionErr.ErrSessionNotFound
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

	result := &SessionResponse{
		ID:           session.ID,
		Name:         session.Name,
		IsActive:     session.IsActive,
		Members:      s.buildSessionMembers(members),
		Transactions: s.buildSessionTransactions(lines),
		Total:        session.Total,
	}

	return result, nil
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
