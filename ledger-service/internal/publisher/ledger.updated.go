package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const TopicLedgerUpdated = "ledger.updated"

type LedgerUpdatedEvent struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	SessionID     uuid.UUID `json:"session_id"`
	UserID        uuid.UUID `json:"user_id"`
	Amount        float64   `json:"amount"`
}

type LedgerUpdatedPublisher interface {
	Publish(ctx context.Context, ev LedgerUpdatedEvent) error
}

type LedgerUpdatedKafkaPublisher struct {
	writer *kafka.Writer
}

func NewLedgerUpdatedKafkaPublisher(brokers []string) LedgerUpdatedPublisher {
	return &LedgerUpdatedKafkaPublisher{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: TopicLedgerUpdated,
		},
	}
}

func (p *LedgerUpdatedKafkaPublisher) Publish(ctx context.Context, ev LedgerUpdatedEvent) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("failed to marshal ledger updated event %v: %w", ev.ID, err)
	}

	msg := kafka.Message{
		Key:   []byte(ev.TransactionID.String()),
		Value: data,
		Headers: []kafka.Header{
			{"version", []byte("1.0.0")},
			{"source", []byte("ledger-service")},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.writer.Topic, err)
	}

	return nil
}
