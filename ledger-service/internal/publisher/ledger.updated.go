package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
)

const TopicLedgerUpdated = "ledger.updated"

type LedgerUpdatedPublisher interface {
	Publish(ctx context.Context, ev event.LedgerUpdateEvent) error
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

func (p *LedgerUpdatedKafkaPublisher) Publish(ctx context.Context, ev event.LedgerUpdateEvent) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("failed to marshal ledger updated event %v: %w", ev.ID, err)
	}

	msg := kafka.Message{
		Key:   []byte(ev.TransactionID),
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
