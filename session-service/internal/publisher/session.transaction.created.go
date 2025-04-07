package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/session-service/internal/events"
)

type SessionTransactionCreatedPublisher interface {
	Publish(ctx context.Context, ev events.SessionTransactionCreatedEvent) error
}

type SessionTransactionCreatedKafkaPublisher struct {
	writer *kafka.Writer
}

func NewSessionTransactionCreatedKafkaPublisher(brokers []string) SessionTransactionCreatedPublisher {
	return &SessionTransactionCreatedKafkaPublisher{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: TopicSessionTransactionCreated,
		},
	}
}

func (p *SessionTransactionCreatedKafkaPublisher) Publish(ctx context.Context, ev events.SessionTransactionCreatedEvent) error {
	data, err := json.Marshal(&ev)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(ev.SessionID),
		Value: data,
		Headers: []kafka.Header{
			{"version", []byte("1.0.0")},
			{"source", []byte("session-service")},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.writer.Topic, err)
	}

	return nil
}
