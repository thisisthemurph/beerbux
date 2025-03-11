package fake

import (
	"context"
	"github.com/thisisthemurph/beerbux/transaction-service/internal/publisher"
)

type TransactionCreatedSuccessPublisher struct{}

func NewFakeTransactionCreatedPublisher() publisher.TransactionCreatedPublisher {
	return &TransactionCreatedSuccessPublisher{}
}

func (p *TransactionCreatedSuccessPublisher) Publish(ctx context.Context, event publisher.TransactionCreatedEvent) error {
	return nil
}
