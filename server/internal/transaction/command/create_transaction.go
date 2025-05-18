package command

import (
	"beerbux/internal/common/history"
	"beerbux/internal/transaction/db"
	"beerbux/pkg/dbtx"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/thisisthemurph/fn"
)

const (
	MinTransactionAmount = 0.5
	MaxTransactionAmount = 2.0
)

var (
	ErrCreatorCannotBeMember      = errors.New("creator cannot be a member of the transaction")
	ErrInactiveSession            = errors.New("session is inactive")
	ErrNotAllMembersPartOfSession = errors.New("member is not part of the session")
	ErrMemberAmountRequired       = errors.New("member amount is required")
	ErrMemberAmountTooLow         = errors.New("member amount must be at least 0.5")
	ErrMemberAmountTooHigh        = errors.New("member amount cannot be more than 2")
	ErrSessionNotFound            = errors.New("session not found")
)

type CreateTransactionCommand struct {
	dbtx.TX
	Queries              *db.Queries
	SessionHistoryWriter history.SessionHistoryWriter
}

func NewCreateTransactionCommand(tx dbtx.TX, queries *db.Queries, historyWriter history.SessionHistoryWriter) *CreateTransactionCommand {
	return &CreateTransactionCommand{
		TX:                   tx,
		Queries:              queries,
		SessionHistoryWriter: historyWriter,
	}
}

type TransactionLine struct {
	MemberID uuid.UUID `json:"userId"`
	Amount   float64   `json:"amount"`
}

type CreateTransactionRequest struct {
	SessionID uuid.UUID         `json:"sessionId"`
	CreatorID uuid.UUID         `json:"creatorId"`
	Lines     []TransactionLine `json:"amounts"`
}

type TransactionResponse struct {
	ID uuid.UUID `json:"id"`
}

func (cmd *CreateTransactionCommand) Execute(ctx context.Context, r CreateTransactionRequest) (*TransactionResponse, error) {
	var err error
	if err = cmd.validateRequest(ctx, r); err != nil {
		return nil, err
	}

	tx, err := cmd.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := cmd.Queries.WithTx(tx)

	transaction, err := qtx.CreateTransaction(ctx, db.CreateTransactionParams{
		SessionID: r.SessionID,
		MemberID:  r.CreatorID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	for _, memberAmount := range r.Lines {
		_, err = qtx.CreateTransactionLine(ctx, db.CreateTransactionLineParams{
			TransactionID: transaction.ID,
			MemberID:      memberAmount.MemberID,
			Amount:        memberAmount.Amount,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create line for transaction: %w", err)
		}
		err = qtx.CreateLedgerEntry(ctx, db.CreateLedgerEntryParams{
			TransactionID: transaction.ID,
			UserID:        memberAmount.MemberID,
			Amount:        memberAmount.Amount,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create credit ledger entry: %w", err)
		}
		err = qtx.CreateLedgerEntry(ctx, db.CreateLedgerEntryParams{
			TransactionID: transaction.ID,
			UserID:        r.CreatorID,
			Amount:        -memberAmount.Amount,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create debit ledger entry: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if historyErr := cmd.SessionHistoryWriter.CreateTransactionCreatedEvent(ctx, r.SessionID, r.CreatorID, history.TransactionHistory{
		TransactionID: transaction.ID,
		Lines: fn.Map(r.Lines, func(tl TransactionLine) history.TransactionHistoryLine {
			return history.TransactionHistoryLine{
				MemberID: tl.MemberID,
				Amount:   tl.Amount,
			}
		}),
	}); historyErr != nil {
		return nil, historyErr
	}

	return &TransactionResponse{
		ID: transaction.ID,
	}, nil
}

func (cmd *CreateTransactionCommand) validateRequest(ctx context.Context, r CreateTransactionRequest) error {
	if r.Lines == nil || len(r.Lines) == 0 {
		return ErrMemberAmountRequired
	}

	for _, m := range r.Lines {
		if m.Amount < MinTransactionAmount {
			return ErrMemberAmountTooLow
		}
		if m.Amount > MaxTransactionAmount {
			return ErrMemberAmountTooHigh
		}
	}

	memberLookup := fn.Map(r.Lines, func(ma TransactionLine) uuid.UUID {
		return ma.MemberID
	})

	if fn.Contains(memberLookup, r.CreatorID) {
		return ErrCreatorCannotBeMember
	}

	return cmd.validateSession(ctx, r.SessionID, memberLookup)
}

func (cmd *CreateTransactionCommand) validateSession(ctx context.Context, sessionID uuid.UUID, memberLookup []uuid.UUID) error {
	session, err := cmd.Queries.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSessionNotFound
		}
		return fmt.Errorf("failed getting session: %w", err)
	}

	if !session.IsActive {
		return ErrInactiveSession
	}

	sessionMemberIDs, err := cmd.Queries.GetSessionMemberIDs(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed getting session member IDs: %w", err)
	}

	if len(sessionMemberIDs) < len(memberLookup) {
		return ErrNotAllMembersPartOfSession
	}

	for _, mid := range memberLookup {
		if !fn.Contains(sessionMemberIDs, mid) {
			return ErrNotAllMembersPartOfSession
		}
	}

	return nil
}
