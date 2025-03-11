package fake

import (
	"context"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/publisher"
)

type LedgerUpdatedSuccessPublisher struct{}

func NewFakeLedgerUpdatedPublisher() publisher.LedgerUpdatedPublisher {
	return &LedgerUpdatedSuccessPublisher{}
}

func (p *LedgerUpdatedSuccessPublisher) Publish(ctx context.Context, ev publisher.LedgerUpdatedEvent) error {
	return nil
}
