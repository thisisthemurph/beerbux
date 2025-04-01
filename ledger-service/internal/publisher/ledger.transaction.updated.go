package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
)

const TopicLedgerTransactionUpdated = "ledger.transaction.updated"

type LedgerTransactionUpdatedPublisher interface {
	Publish(ctx context.Context, ledgerItems []event.LedgerUpdateEvent) error
}

type LedgerTransactionUpdatedKafkaPublisher struct {
	writer *kafka.Writer
}

func NewLedgerTransactionUpdatedKafkaPublisher(brokers []string) LedgerTransactionUpdatedPublisher {
	return &LedgerTransactionUpdatedKafkaPublisher{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(brokers...),
			Topic: TopicLedgerTransactionUpdated,
		},
	}
}

func (p *LedgerTransactionUpdatedKafkaPublisher) Publish(ctx context.Context, ledgerItems []event.LedgerUpdateEvent) error {
	transactionEvent := event.LedgerTransactionUpdatedEvent{}

	for i, item := range ledgerItems {
		if i == 0 {
			transactionEvent.TransactionID = item.TransactionID
			transactionEvent.SessionID = item.SessionID
			transactionEvent.UserID = item.UserID
			transactionEvent.Amounts = make([]event.LedgerUpdateCompleteMemberAmount, 0, len(ledgerItems))
		}

		// Add the amount to the total if the user is not the creator.
		if item.UserID != transactionEvent.UserID {
			transactionEvent.Total += item.Amount * -1 // The amount should be positive here as we are indicating transaction value.
			transactionEvent.Amounts = append(transactionEvent.Amounts, event.LedgerUpdateCompleteMemberAmount{
				UserID: item.UserID,
				Amount: item.Amount * -1, // The amount should be positive here as we are indicating transaction value.
			})
		}
	}

	data, err := json.Marshal(transactionEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal %v event %v: %w", p.writer.Topic, transactionEvent.TransactionID, err)
	}

	msg := kafka.Message{
		Key:   []byte(transactionEvent.TransactionID),
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
