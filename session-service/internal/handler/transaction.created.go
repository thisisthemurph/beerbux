package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/session-service/internal/events"
	"github.com/thisisthemurph/beerbux/session-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/fn"
)

type TransactionCreatedMessageHandler struct {
	sessionRepository           *repository.SessionQueriesWrapper
	transactionCreatedPublisher publisher.SessionTransactionCreatedPublisher
}

func NewTransactionCreatedMessageHandler(
	sessionRepository *repository.SessionQueriesWrapper,
	p publisher.SessionTransactionCreatedPublisher) KafkaMessageHandler {
	return &TransactionCreatedMessageHandler{
		sessionRepository:           sessionRepository,
		transactionCreatedPublisher: p,
	}
}

type MemberAmount struct {
	MemberID string  `json:"user_id"`
	Amount   float64 `json:"amount"`
}

type TransactionCreatedEvent struct {
	TransactionID string         `json:"transaction_id"`
	SessionID     string         `json:"session_id"`
	CreatorID     string         `json:"creator_id"`
	Amounts       []MemberAmount `json:"member_amounts"`
}

func (h TransactionCreatedMessageHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var event TransactionCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	tx, err := h.sessionRepository.Transaction.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	qtx := h.sessionRepository.WithTx(tx)

	t, err := qtx.AddTransaction(ctx, session.AddTransactionParams{
		ID:        event.TransactionID,
		SessionID: event.SessionID,
		MemberID:  event.CreatorID,
		CreatedAt: msg.Time,
	})
	if err != nil {
		return fmt.Errorf("failed adding transaction: %w", err)
	}

	total := 0.0
	for _, a := range event.Amounts {
		if a.Amount > 0 {
			total += a.Amount
		}
		_, err := qtx.AddTransactionLine(ctx, session.AddTransactionLineParams{
			TransactionID: t.ID,
			MemberID:      a.MemberID,
			Amount:        a.Amount,
		})
		if err != nil {
			return fmt.Errorf("failed adding line to transaction: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed committing transaction: %w", err)
	}

	_ = h.transactionCreatedPublisher.Publish(ctx, events.SessionTransactionCreatedEvent{
		SessionID:     event.SessionID,
		TransactionID: event.TransactionID,
		CreatorID:     event.CreatorID,
		Total:         total,
		Members: fn.Map(event.Amounts, func(ma MemberAmount) events.SessionMember {
			return events.SessionMember{
				ID: ma.MemberID,
			}
		}),
	})

	return nil
}
