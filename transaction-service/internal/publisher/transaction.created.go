package publisher

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type TransactionCreatedMemberAmount struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type TransactionCreatedEvent struct {
	TransactionID string                           `json:"transaction_id"`
	CreatorID     string                           `json:"creator_id"`
	SessionID     string                           `json:"session_id"`
	MemberAmounts []TransactionCreatedMemberAmount `json:"member_amounts"`
}

type TransactionCreatedPublisher interface {
	Publish(ctx context.Context, t TransactionCreatedEvent) error
}

type TransactionCreatedKafkaPublisher struct {
	writer *kafka.Writer
}

func NewTransactionCreatedKafkaPublisher(brokers []string) TransactionCreatedPublisher {
	return &TransactionCreatedKafkaPublisher{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: TopicTransactionCreated,
		},
	}
}

func (p *TransactionCreatedKafkaPublisher) Publish(ctx context.Context, ev TransactionCreatedEvent) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: data,
		Headers: []kafka.Header{
			{"version", []byte("1.0.0")},
			{"source", []byte("transaction-service")},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return err
	}

	return nil
}
