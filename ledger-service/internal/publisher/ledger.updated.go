package publisher

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
)

const SubjectLedgerUpdated = "ledger.updated"

type LedgerUpdatedEventData struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	SessionID     uuid.UUID `json:"session_id"`
	UserID        uuid.UUID `json:"user_id"`
	Amount        float64   `json:"amount"`
}

type LedgerUpdatedEvent struct {
	event.Metadata
	Data LedgerUpdatedEventData `json:"data"`
}

type LedgerUpdatedPublisher interface {
	Publish(id, transactionID, sessionID, userID uuid.UUID, amount float64) error
}

type LedgerUpdatedNatsPublisher struct {
	nc      *nats.Conn
	subject string
}

func NewLedgerUpdatedNatsPublisher(nc *nats.Conn) LedgerUpdatedPublisher {
	return &LedgerUpdatedNatsPublisher{
		nc:      nc,
		subject: SubjectLedgerUpdated,
	}
}

func (p *LedgerUpdatedNatsPublisher) Publish(id, transactionID, sessionID, userID uuid.UUID, amount float64) error {
	msg := LedgerUpdatedEvent{
		Metadata: event.NewMetadata(SubjectLedgerUpdated, "1.0.0", transactionID.String()),
		Data: LedgerUpdatedEventData{
			ID:            id,
			TransactionID: transactionID,
			SessionID:     sessionID,
			UserID:        userID,
			Amount:        amount,
		},
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := p.nc.Publish(p.subject, msgData); err != nil {
		return err
	}

	return nil
}
