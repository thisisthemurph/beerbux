package fake

import (
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/publisher"
)

type LedgerUpdatedSuccessPublisher struct{}

func NewFakeLedgerUpdatedPublisher() publisher.LedgerUpdatedPublisher {
	return &LedgerUpdatedSuccessPublisher{}
}

func (p *LedgerUpdatedSuccessPublisher) Publish(id, transactionID, sessionID, userID uuid.UUID, amount float64) error {
	return nil
}
