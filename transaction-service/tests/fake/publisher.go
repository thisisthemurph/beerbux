package fake

import "github.com/thisisthemurph/beerbux/transaction-service/internal/publisher"

type TransactionCreatedSuccessPublisher struct{}

func NewFakeTransactionCreatedPublisher() publisher.TransactionCreatedPublisher {
	return &TransactionCreatedSuccessPublisher{}
}

func (p *TransactionCreatedSuccessPublisher) Publish(event publisher.TransactionCreatedEventData) error {
	return nil
}
