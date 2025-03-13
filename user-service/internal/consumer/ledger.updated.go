package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/ledger"
)

type LedgerUpdatedEvent struct {
	ID            string  `json:"id"`
	TransactionID string  `json:"transaction_id"`
	SessionID     string  `json:"session_id"`
	UserID        string  `json:"user_id"`
	ParticipantID string  `json:"participant_id"`
	Amount        float64 `json:"amount"`
}

type LedgerUpdatedKafkaConsumer struct {
	Reader               KafkaReader
	Logger               *slog.Logger
	UserLedgerRepository *ledger.Queries
}

func NewLedgerUpdatedKafkaConsumer(logger *slog.Logger, brokers []string, topic string, repo *ledger.Queries) *LedgerUpdatedKafkaConsumer {
	return &LedgerUpdatedKafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: "user-service",
		}),
		Logger:               logger,
		UserLedgerRepository: repo,
	}
}

func (c *LedgerUpdatedKafkaConsumer) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.Logger.Debug("Kafka consumer shutting down")
			return
		default:
			msg, err := c.Reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil || errors.Is(err, context.Canceled) {
					return
				}
				c.Logger.Error("Failed to read message", "error", err)
				continue
			}

			var ev LedgerUpdatedEvent
			if err := json.Unmarshal(msg.Value, &ev); err != nil {
				c.Logger.Error("Failed to unmarshal message", "error", err, "offset", msg.Offset)
				continue
			}

			err = c.UserLedgerRepository.InsertUserLedger(ctx, ledger.InsertUserLedgerParams{
				UserID:        ev.UserID,
				ParticipantID: ev.ParticipantID,
				Amount:        ev.Amount,
			})

			if err != nil {
				c.Logger.Error("Failed to insert user ledger", "error", err, "user_id", ev.UserID, "offset", msg.Offset)
				continue
			}

			c.Logger.Info("Successfully processed message", "offset", msg.Offset, "user_id", ev.UserID)
		}
	}
}

func (c *LedgerUpdatedKafkaConsumer) Close() error {
	return c.Reader.Close()
}
