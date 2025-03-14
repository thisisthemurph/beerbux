package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type UserRegisteredEvent struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserRegisteredProducer interface {
	Publish(ctx context.Context, ev UserRegisteredEvent) error
}

type UserRegisteredKafkaProducer struct {
	writer *kafka.Writer
}

func NewUserRegisteredKafkaProducer(brokers []string) UserRegisteredProducer {
	return &UserRegisteredKafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        TopicUserRegistered,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *UserRegisteredKafkaProducer) Publish(ctx context.Context, ev UserRegisteredEvent) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("failed to marshal user registered event %v: %w", ev.UserID, err)
	}

	msg := kafka.Message{
		Key:   []byte(ev.UserID),
		Value: data,
		Headers: []kafka.Header{
			{"version", []byte("1.0.0")},
			{"source", []byte("auth-service")},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.writer.Topic, err)
	}

	return nil
}
