package queries

import (
	"beerbux/internal/session/db"
	sessionErr "beerbux/internal/session/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/thisisthemurph/fn"
	"time"
)

type GetSessionQuery struct {
	Queries *db.Queries
}

func NewGetSessionQuery(queries *db.Queries) *GetSessionQuery {
	return &GetSessionQuery{
		Queries: queries,
	}
}

type SessionMember struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"isAdmin"`
	IsDeleted bool      `json:"isDeleted"`
}

type SessionTransactionLine struct {
	UserID uuid.UUID `json:"userId"`
	Amount float64   `json:"amount"`
}

type SessionTransaction struct {
	ID        uuid.UUID                `json:"id"`
	UserID    uuid.UUID                `json:"userId"`
	Total     float64                  `json:"total"`
	Lines     []SessionTransactionLine `json:"lines"`
	CreatedAt time.Time                `json:"createdAt"`
}

type SessionResponse struct {
	ID           uuid.UUID            `json:"id"`
	Name         string               `json:"name"`
	IsActive     bool                 `json:"isActive"`
	Members      []SessionMember      `json:"members"`
	Transactions []SessionTransaction `json:"transactions"`
	Total        float64              `json:"total"`
}

func (q *GetSessionQuery) Execute(ctx context.Context, sessionID uuid.UUID) (*SessionResponse, error) {
	s, err := q.Queries.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sessionErr.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session %s: %w", sessionID, err)
	}

	lines, err := q.Queries.GetSessionTransactionLines(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction lines for session %s: %w", sessionID, err)
	}

	members, err := q.Queries.ListSessionMembers(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list session %s members: %w", sessionID, err)
	}

	result := &SessionResponse{
		ID:           s.ID,
		Name:         s.Name,
		IsActive:     s.IsActive,
		Members:      q.BuildSessionMembers(members),
		Transactions: q.BuildSessionTransactions(lines),
		Total:        s.Total,
	}

	return result, nil
}

func (q *GetSessionQuery) BuildSessionMembers(members []db.ListSessionMembersRow) []SessionMember {
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

func (q *GetSessionQuery) BuildSessionTransactions(lines []db.GetSessionTransactionLinesRow) []SessionTransaction {
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
