package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/pkg/nullish"
)

type UserCreatedEvent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

type UserCreatedPublisher interface {
	Publish(u user.User) error
}

type UserCreatedKafkaPublisher struct {
	writer *kafka.Writer
}

func NewUserCreatedKafkaPublisher(brokers []string) UserCreatedPublisher {
	return &UserCreatedKafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        TopicUserCreated,
			Balancer:     nil,
			BatchSize:    0,
			BatchTimeout: 10 * time.Millisecond,
			Async:        false,
		},
	}
}

func (p *UserCreatedKafkaPublisher) Publish(u user.User) error {
	ev := UserCreatedEvent{
		ID:       u.ID,
		Name:     u.Name,
		Username: u.Username,
		Bio:      nullish.StringOrEmpty(u.Bio),
	}

	data, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("failed to marshal user %v: %w", u.ID, err)
	}

	msg := kafka.Message{
		Value:   data,
		Headers: makeKafkaHeaders(u),
	}

	if err := p.writer.WriteMessages(context.TODO(), msg); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.writer.Topic, err)
	}

	return nil
}

func makeKafkaHeaders(u user.User) []kafka.Header {
	return []kafka.Header{
		{"version", []byte("1.0.0")},
		{"source", []byte("user-service")},
		{"correlation_id", []byte(u.ID)},
	}
}
