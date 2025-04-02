package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	
	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
)

const TopicLedgerUserTotalsCalculated = "ledger.user.totals.calculated"

type LedgerUserTotalsCalculatedPublisher interface {
	Publish(ctx context.Context, ev event.UserTotalsEvent) error
}

type LedgerUserTotalsCalculatedKafkaPublisher struct {
	writer *kafka.Writer
}

func NewLedgerUserTotalsCalculatedKafkaPublisher(brokers []string) LedgerUserTotalsCalculatedPublisher {
	return &LedgerUserTotalsCalculatedKafkaPublisher{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: TopicLedgerUserTotalsCalculated,
		},
	}
}

func (p LedgerUserTotalsCalculatedKafkaPublisher) Publish(ctx context.Context, ev event.UserTotalsEvent) error {
	data, err := json.Marshal(&ev)
	if err != nil {
		return fmt.Errorf("%s: failed marshalling event: %w", TopicLedgerUserTotalsCalculated, err)
	}

	msg := kafka.Message{
		Key:   []byte(ev.UserID),
		Value: data,
		Headers: []kafka.Header{
			{"version", []byte("1.0.0")},
			{"source", []byte("ledger.service")},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("%s: failed to publish message: %w", TopicLedgerUserTotalsCalculated, err)
	}

	return nil
}
