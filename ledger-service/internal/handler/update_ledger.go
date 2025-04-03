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
	ledgerUpdatedChan chan<- LedgerTransaction
	logger            *slog.Logger
}

func NewUpdateLedgerHandler(
	ledgerRepository *repository.LedgerQueriesWrapper,
	ledgerUpdatedChan chan<- LedgerTransaction,
	logger *slog.Logger,
) *UpdateLedgerHandler {
	return &UpdateLedgerHandler{
		ledgerRepository:  ledgerRepository,
		ledgerUpdatedChan: ledgerUpdatedChan,
		logger:            logger,
	}
}

type LedgerTransaction struct {
	CreatorID     string       `json:"creator_id"`
	TransactionID string       `json:"transaction_id"`
	SessionID     string       `json:"session_id"`
	LedgerItems   []LedgerItem `json:"ledger_items"`
}

// GetAllMemberIDs returns a slice of all member IDs involved in the transaction, including the creator.
func (lt LedgerTransaction) GetAllMemberIDs() []string {
	memberIDs := make([]string, 0, len(lt.LedgerItems))

	memberIDs = append(memberIDs, lt.CreatorID)
	for _, item := range lt.LedgerItems {
		if item.ParticipantID != lt.CreatorID {
			memberIDs = append(memberIDs, item.ParticipantID)
		}
	}
	return memberIDs
}

type LedgerItem struct {
	UserID        string  `json:"user_id"`
	ParticipantID string  `json:"participant_id"`
	Amount        float64 `json:"amount"`
	Type          string  `json:"type"`
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

	ledgerTransaction := h.buildLedgerTransactionFromTransactionCreatedEvent(ev)
	err := h.updateLedger(ctx, ledgerTransaction)
	if err != nil {
		h.logger.Error("failed to update ledger", "transactionID", ev.TransactionID, "error", err)
		return err
	}

	h.ledgerUpdatedChan <- ledgerTransaction
	return nil
}

func (h *UpdateLedgerHandler) buildLedgerTransactionFromTransactionCreatedEvent(ev event.TransactionCreatedEvent) LedgerTransaction {
	ledgerTransaction := LedgerTransaction{
		CreatorID:     ev.CreatorID,
		TransactionID: ev.TransactionID,
		SessionID:     ev.SessionID,
		LedgerItems:   make([]LedgerItem, 0, len(ev.MemberAmounts)*2),
	}

	for _, ma := range ev.MemberAmounts {
		ledgerTransaction.LedgerItems = append(ledgerTransaction.LedgerItems, LedgerItem{
			UserID:        ma.UserID,
			ParticipantID: ev.CreatorID,
			Amount:        ma.Amount,
			Type:          "credit",
		})

		ledgerTransaction.LedgerItems = append(ledgerTransaction.LedgerItems, LedgerItem{
			UserID:        ev.CreatorID,
			ParticipantID: ma.UserID,
			Amount:        -ma.Amount,
			Type:          "debit",
		})
	}

	return ledgerTransaction
}

func (h *UpdateLedgerHandler) updateLedger(ctx context.Context, ledgerTransaction LedgerTransaction) error {
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

	for _, li := range ledgerTransaction.LedgerItems {
		params := ledger.InsertLedgerParams{
			ID:            uuid.NewString(),
			TransactionID: ledgerTransaction.TransactionID,
			SessionID:     ledgerTransaction.SessionID,
			UserID:        li.UserID,
			Amount:        li.Amount,
		}
		if err = qtx.InsertLedger(ctx, params); err != nil {
			return fmt.Errorf("failed to insert ledger for transactionID %q: %w", ledgerTransaction.TransactionID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction for transactionID: %w", err)
	}

	return nil
}
