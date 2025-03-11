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
	"github.com/thisisthemurph/beerbux/ledger-service/pkg/fn"
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

type LedgerUpdateResult struct {
	ID            uuid.UUID
	TransactionID uuid.UUID
	SessionID     uuid.UUID
	UserID        uuid.UUID
	Amount        float64
}

func (h *UpdateLedgerHandler) Handle(ctx context.Context, ev event.TransactionCreatedEvent) ([]*LedgerUpdateResult, error) {
	if len(ev.MemberAmounts) == 0 {
		h.logger.Error("no member amounts provided", "transactionID", ev.TransactionID)
		return nil, fmt.Errorf("no member amounts provided for transactionID %q", ev.TransactionID)
	}

	insertedLedgerItems, err := h.updateLedger(ctx, ev)
	if err != nil {
		h.logger.Error("failed to update ledger", "transactionID", ev.TransactionID, "error", err)
		return nil, err
	}

	results := fn.Map(insertedLedgerItems, func(l ledger.InsertLedgerParams) *LedgerUpdateResult {
		return &LedgerUpdateResult{
			ID:            uuid.MustParse(l.ID),
			TransactionID: ev.TransactionID,
			SessionID:     ev.SessionID,
			UserID:        uuid.MustParse(l.UserID),
			Amount:        l.Amount,
		}
	})

	return results, nil
}

func (h *UpdateLedgerHandler) updateLedger(ctx context.Context, data event.TransactionCreatedEvent) ([]ledger.InsertLedgerParams, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := h.ledgerRepository.Transaction.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction for transactionID %q: %w", data.TransactionID.String(), err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	qtx := h.ledgerRepository.WithTx(tx)

	var inserts []ledger.InsertLedgerParams
	for _, member := range data.MemberAmounts {
		// Credit entries for the members
		inserts = append(inserts, ledger.InsertLedgerParams{
			ID:            uuid.NewString(),
			TransactionID: data.TransactionID.String(),
			SessionID:     data.SessionID.String(),
			UserID:        member.UserID.String(),
			Amount:        member.Amount,
		})

		// Debit entries for the creator
		inserts = append(inserts, ledger.InsertLedgerParams{
			ID:            uuid.NewString(),
			TransactionID: data.TransactionID.String(),
			SessionID:     data.SessionID.String(),
			UserID:        data.CreatorID.String(),
			Amount:        -member.Amount,
		})
	}

	for _, ledgerRecord := range inserts {
		if err = qtx.InsertLedger(ctx, ledgerRecord); err != nil {
			return nil, fmt.Errorf("failed to insert ledger for transactionID %q: %w", data.TransactionID.String(), err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction for transactionID %q: %w", data.TransactionID.String(), err)
	}

	return inserts, nil
}
