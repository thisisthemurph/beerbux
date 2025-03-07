package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository/ledger"
	"github.com/thisisthemurph/beerbux/ledger-service/pkg/semanticversion"
)

type TransactionCreatedMsgHandler struct {
	ledgerRepository *repository.LedgerQueriesWrapper
	logger           *slog.Logger
}

func NewTransactionCreatedMsgHandler(ledgerRepository *repository.LedgerQueriesWrapper, logger *slog.Logger) *TransactionCreatedMsgHandler {
	return &TransactionCreatedMsgHandler{
		ledgerRepository: ledgerRepository,
		logger:           logger,
	}
}

func (h *TransactionCreatedMsgHandler) Handle(msg *nats.Msg) {
	ctx := context.Background()

	var ev TransactionCreatedEvent
	err := json.Unmarshal(msg.Data, &ev)
	if err != nil {
		h.logger.Error("failed to unmarshal event", "error", err)
		return
	}

	v, err := semanticversion.Parse(ev.Version)
	if err != nil {
		h.logger.Error("failed to parse event version", "error", err, "version", ev.Version)
		return
	}
	if v.Major != 1 {
		h.logger.Error("unexpected event major version, expected major version 1", "version", ev.Version)
		return
	}

	err = h.updateLedger(ctx, ev.Data)
	if err != nil {
		h.logger.Error("failed to update ledger", "transactionID", ev.Data.TransactionID, "error", err)
	}
}

func (h *TransactionCreatedMsgHandler) updateLedger(ctx context.Context, data TransactionCreatedEventData) error {
	var err error
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := h.ledgerRepository.Transaction.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	qtx := h.ledgerRepository.WithTx(tx)

	for _, m := range data.MemberAmounts {
		err = qtx.InsertLedger(ctx, ledger.InsertLedgerParams{
			ID:            uuid.NewString(),
			TransactionID: data.TransactionID.String(),
			SessionID:     data.SessionID.String(),
			UserID:        m.UserID.String(),
			Amount:        m.Amount,
		})
		if err != nil {
			return fmt.Errorf("failed to insert ledger: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
