package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"time"
)

type UserUpdatedEvent struct {
	UserID        string                 `json:"user_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"`
}

type UserUpdatedPublisher interface {
	Publish(original user.User, new user.User) error
}

type UserUpdatedKafkaPublisher struct {
	writer *kafka.Writer
}

func NewUserUpdatedKafkaPublisher(brokers []string) UserUpdatedPublisher {
	return &UserUpdatedKafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        TopicUserUpdated,
			Balancer:     nil,
			BatchSize:    0,
			BatchTimeout: 10 * time.Millisecond,
			Async:        false,
		},
	}
}

func (p *UserUpdatedKafkaPublisher) Publish(original, new user.User) error {
	ev := UserUpdatedEvent{
		UserID:        original.ID,
		UpdatedFields: p.determineUpdatedFields(original, new),
	}

	data, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("failed to marshal user %v: %w", new.ID, err)
	}

	msg := kafka.Message{
		Value:   data,
		Key:     []byte(original.ID),
		Headers: makeKafkaHeaders(new),
	}

	if err := p.writer.WriteMessages(context.TODO(), msg); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.writer.Topic, err)
	}

	return nil
}

func (p *UserUpdatedKafkaPublisher) determineUpdatedFields(original, new user.User) map[string]interface{} {
	updatedFields := make(map[string]interface{})

	if original.Name != new.Name {
		updatedFields["name"] = new.Name
	}

	if original.Username != new.Username {
		updatedFields["username"] = new.Username
	}

	if original.Bio != new.Bio {
		updatedFields["bio"] = new.Bio
	}

	updatedFields["updated_at"] = new.UpdatedAt

	return updatedFields
}
