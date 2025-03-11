package handler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository/ledger"
)

type UpdateLedgerHandler struct {
	ledgerRepository *repository.LedgerQueriesWrapper
	logger           *slog.Logger
}

func NewUpdateLedgerHandler(
	ledgerRepository *repository.LedgerQueriesWrapper,
	logger *slog.Logger,
) *UpdateLedgerHandler {
	return &UpdateLedgerHandler{
		ledgerRepository: ledgerRepository,
		logger:           logger,
	}
}

// MemberTransaction is a snapshot of an individual transaction within an event.TransactionCreatedEvent.
// The transaction will be between the creator and one of the members or one of the members and the creator.
type MemberTransaction struct {
	ledger.InsertLedgerParams
	ParticipantID string
}

func (h *UpdateLedgerHandler) Handle(ctx context.Context, ev event.TransactionCreatedEvent) ([]MemberTransaction, error) {
	if len(ev.MemberAmounts) == 0 {
		h.logger.Error("no member amounts provided", "transactionID", ev.TransactionID)
		return nil, fmt.Errorf("no member amounts provided for transactionID %q", ev.TransactionID)
	}

	transactions := h.getTransactionsFromEvent(ev)
	err := h.updateLedger(ctx, transactions)
	if err != nil {
		h.logger.Error("failed to update ledger", "transactionID", ev.TransactionID, "error", err)
		return nil, err
	}

	return transactions, nil
}

// getTransactionsFromEvent generates a list of MemberTransaction from a given event.TransactionCreatedEvent.
// A MemberTransaction is generated for each Creator and Member combination as well as the reverse.
func (h *UpdateLedgerHandler) getTransactionsFromEvent(ev event.TransactionCreatedEvent) []MemberTransaction {
	transactions := make([]MemberTransaction, 0, len(ev.MemberAmounts)*2)
	for _, ma := range ev.MemberAmounts {
		// Credit entries for the members
		transactions = append(transactions, MemberTransaction{
			InsertLedgerParams: ledger.InsertLedgerParams{
				ID:            uuid.NewString(),
				TransactionID: ev.TransactionID,
				SessionID:     ev.SessionID,
				UserID:        ma.UserID,
				Amount:        ma.Amount,
			},
			ParticipantID: ev.CreatorID,
		})

		// Debit entries for the creator
		transactions = append(transactions, MemberTransaction{
			InsertLedgerParams: ledger.InsertLedgerParams{
				ID:            uuid.NewString(),
				TransactionID: ev.TransactionID,
				SessionID:     ev.SessionID,
				UserID:        ev.CreatorID,
				Amount:        -ma.Amount,
			},
			ParticipantID: ma.UserID,
		})
	}

	return transactions
}

func (h *UpdateLedgerHandler) updateLedger(ctx context.Context, transactions []MemberTransaction) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := h.ledgerRepository.Transaction.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for transactionID: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	qtx := h.ledgerRepository.WithTx(tx)

	for _, mt := range transactions {
		if err = qtx.InsertLedger(ctx, mt.InsertLedgerParams); err != nil {
			return fmt.Errorf("failed to insert ledger for transactionID %q: %w", mt.InsertLedgerParams.TransactionID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction for transactionID: %w", err)
	}

	return nil
}
