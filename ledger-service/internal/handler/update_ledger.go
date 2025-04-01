package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository/ledger"
)

type UpdateLedgerHandler struct {
	ledgerRepository  *repository.LedgerQueriesWrapper
	ledgerUpdatedChan chan<- []event.LedgerUpdateEvent
	logger            *slog.Logger
}

func NewUpdateLedgerHandler(
	ledgerRepository *repository.LedgerQueriesWrapper,
	ledgerUpdatedChan chan<- []event.LedgerUpdateEvent,
	logger *slog.Logger,
) *UpdateLedgerHandler {
	return &UpdateLedgerHandler{
		ledgerRepository:  ledgerRepository,
		ledgerUpdatedChan: ledgerUpdatedChan,
		logger:            logger,
	}
}

// Handle processes the incoming Kafka message and updates the ledger.
// The transaction.created event is converted into individual transactions for each member/creator combination.
func (h *UpdateLedgerHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var ev event.TransactionCreatedEvent
	if err := json.Unmarshal(msg.Value, &ev); err != nil {
		h.logger.Error("failed to unmarshal event", "error", err)
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	if len(ev.MemberAmounts) == 0 {
		h.logger.Error("no member amounts provided", "transactionID", ev.TransactionID)
		return fmt.Errorf("no member amounts provided for transactionID %q", ev.TransactionID)
	}

	transactions := h.getTransactionsFromEvent(ev)
	err := h.updateLedger(ctx, transactions)
	if err != nil {
		h.logger.Error("failed to update ledger", "transactionID", ev.TransactionID, "error", err)
		return err
	}

	h.ledgerUpdatedChan <- transactions
	return nil
}

// getTransactionsFromEvent generates a list of event.LedgerUpdateEvent from a given event.TransactionCreatedEvent.
// An event.LedgerUpdateEvent is generated for each creator/member combination as well as the reverse.
func (h *UpdateLedgerHandler) getTransactionsFromEvent(ev event.TransactionCreatedEvent) []event.LedgerUpdateEvent {
	transactions := make([]event.LedgerUpdateEvent, 0, len(ev.MemberAmounts)*2)
	for _, ma := range ev.MemberAmounts {
		// Credit entries for the members
		transactions = append(transactions, event.LedgerUpdateEvent{
			ID:            uuid.NewString(),
			TransactionID: ev.TransactionID,
			SessionID:     ev.SessionID,
			UserID:        ma.UserID,
			Amount:        ma.Amount,
			ParticipantID: ev.CreatorID,
		})

		// Debit entries for the creator
		transactions = append(transactions, event.LedgerUpdateEvent{
			ID:            uuid.NewString(),
			TransactionID: ev.TransactionID,
			SessionID:     ev.SessionID,
			UserID:        ev.CreatorID,
			Amount:        -ma.Amount,
			ParticipantID: ma.UserID,
		})
	}

	return transactions
}

func (h *UpdateLedgerHandler) updateLedger(ctx context.Context, transactions []event.LedgerUpdateEvent) error {
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
		params := ledger.InsertLedgerParams{
			ID:            mt.ID,
			TransactionID: mt.TransactionID,
			SessionID:     mt.SessionID,
			UserID:        mt.UserID,
			Amount:        mt.Amount,
		}
		if err = qtx.InsertLedger(ctx, params); err != nil {
			return fmt.Errorf("failed to insert ledger for transactionID %q: %w", mt.TransactionID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction for transactionID: %w", err)
	}

	return nil
}
