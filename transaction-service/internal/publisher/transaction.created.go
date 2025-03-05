package publisher

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

const SubjectTransactionCreated = "transaction.created"

type TransactionCreatedMemberAmount struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type TransactionCreatedEventData struct {
	TransactionID string                           `json:"transaction_id"`
	CreatorID     string                           `json:"creator_id"`
	SessionID     string                           `json:"session_id"`
	MemberAmounts []TransactionCreatedMemberAmount `json:"member_amounts"`
}

type TransactionCreatedEvent struct {
	EventMetadata
	Data TransactionCreatedEventData `json:"data"`
}

type TransactionCreatedPublisher interface {
	Publish(t TransactionCreatedEventData) error
}

type TransactionCreatedNatsPublisher struct {
	nc      *nats.Conn
	subject string
}

func NewTransactionCreatedNatsPublisher(nc *nats.Conn) TransactionCreatedPublisher {
	return &TransactionCreatedNatsPublisher{
		nc:      nc,
		subject: SubjectTransactionCreated,
	}
}

func (p *TransactionCreatedNatsPublisher) Publish(ev TransactionCreatedEventData) error {
	msg := TransactionCreatedEvent{
		EventMetadata: NewEventMetadata(SubjectTransactionCreated, "1.0.0", ev.TransactionID),
		Data:          ev,
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal event %v: %w", ev.TransactionID, err)
	}

	if err := p.nc.Publish(p.subject, msgData); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.subject, err)
	}

	return nil
}
