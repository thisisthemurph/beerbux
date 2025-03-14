package fakes

import (
	"context"
	"github.com/thisisthemurph/beerbux/auth-service/internal/producer"
)

type UserRegisteredProducer struct{}

func NewFakeUserRegisteredProducer() producer.UserRegisteredProducer {
	return &UserRegisteredProducer{}
}

func (p *UserRegisteredProducer) Publish(ctx context.Context, ev producer.UserRegisteredEvent) error {
	return nil
}
