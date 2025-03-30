package fake

import (
	"context"
	"github.com/thisisthemurph/beerbux/transaction-service/internal/publisher"
)

type TransactionCreatedSuccessPublisher struct {
	CapturedEvent *publisher.TransactionCreatedEvent
}

func NewFakeTransactionCreatedPublisher() *TransactionCreatedSuccessPublisher {
	return &TransactionCreatedSuccessPublisher{}
}

func (p *TransactionCreatedSuccessPublisher) Publish(ctx context.Context, event publisher.TransactionCreatedEvent) error {
	p.CapturedEvent = &event
	return nil
}
