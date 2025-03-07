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
	"github.com/thisisthemurph/beerbux/ledger-service/pkg/semanticversion"
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

func (h *UpdateLedgerHandler) Handle(ev event.TransactionCreatedEvent) ([]*LedgerUpdateResult, error) {
	ctx := context.Background()

	v, err := semanticversion.Parse(ev.Version)
	if err != nil {
		h.logger.Error("failed to parse event version", "error", err, "version", ev.Version)
		return nil, err
	}
	if v.Major != 1 {
		h.logger.Error("unexpected event major version, expected major version 1", "version", ev.Version)
		return nil, err
	}

	insertedLedgerItems, err := h.updateLedger(ctx, ev.Data)
	if err != nil {
		h.logger.Error("failed to update ledger", "transactionID", ev.Data.TransactionID, "error", err)
		return nil, err
	}

	results := fn.Map(insertedLedgerItems, func(l ledger.InsertLedgerParams) *LedgerUpdateResult {
		return &LedgerUpdateResult{
			ID:            uuid.MustParse(l.ID),
			TransactionID: uuid.MustParse(l.TransactionID),
			SessionID:     uuid.MustParse(l.SessionID),
			UserID:        uuid.MustParse(l.UserID),
			Amount:        l.Amount,
		}
	})

	return results, nil
}

func (h *UpdateLedgerHandler) updateLedger(ctx context.Context, data event.TransactionCreatedEventData) ([]ledger.InsertLedgerParams, error) {
	var err error
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := h.ledgerRepository.Transaction.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	qtx := h.ledgerRepository.WithTx(tx)

	inserts := fn.Map(data.MemberAmounts, func(m event.TransactionCreatedMemberAmount) ledger.InsertLedgerParams {
		return ledger.InsertLedgerParams{
			ID:            uuid.NewString(),
			TransactionID: data.TransactionID.String(),
			SessionID:     data.SessionID.String(),
			UserID:        m.UserID.String(),
			Amount:        m.Amount,
		}
	})

	for _, ledgerRecord := range inserts {
		err = qtx.InsertLedger(ctx, ledgerRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to insert ledger: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return inserts, nil
}
