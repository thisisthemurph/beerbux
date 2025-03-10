package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type SessionMemberAddedPublisher interface {
	Publish(sessionID, memberID string) error
}

type SessionMemberAddedKafkaPublisher struct {
	writer *kafka.Writer
}

func NewSessionMemberAddedKafkaPublisher(brokers []string) SessionMemberAddedPublisher {
	return &SessionMemberAddedKafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        TopicSessionMemberAdded,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

type SessionMemberAddedEventData struct {
	SessionID string `json:"session_id"`
	MemberID  string `json:"member_id"`
}

func (p *SessionMemberAddedKafkaPublisher) Publish(sessionID, memberID string) error {
	data, err := json.Marshal(SessionMemberAddedEventData{
		SessionID: sessionID,
		MemberID:  memberID,
	})

	if err != nil {
		return fmt.Errorf("failed to marshal session member added %v: %w", sessionID, err)
	}

	msg := kafka.Message{
		Key:   []byte(sessionID),
		Value: data,
		Headers: []kafka.Header{
			{"version", []byte("1.0.0")},
			{"source", []byte("session-service")},
		},
	}

	if err := p.writer.WriteMessages(context.TODO(), msg); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.writer.Topic, err)
	}

	return nil
}
