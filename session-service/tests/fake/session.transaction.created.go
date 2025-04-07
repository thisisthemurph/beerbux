package fake

import (
	"context"
	"github.com/thisisthemurph/beerbux/session-service/internal/events"
	"github.com/thisisthemurph/beerbux/session-service/internal/publisher"
)

type SessionTransactionCreatedPublisher struct {
}

func NewFakeSessionTransactionCreatedPublisher() publisher.SessionTransactionCreatedPublisher {
	return &SessionTransactionCreatedPublisher{}
}

func (s SessionTransactionCreatedPublisher) Publish(ctx context.Context, ev events.SessionTransactionCreatedEvent) error {
	return nil
}
