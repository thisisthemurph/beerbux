package fake

import (
	"github.com/thisisthemurph/beerbux/user-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
)

type UserCreatedSuccessPublisher struct{}

func NewFakeUserCreatedPublisher() publisher.UserCreatedPublisher {
	return &UserCreatedSuccessPublisher{}
}

func (p *UserCreatedSuccessPublisher) Publish(u user.User) error {
	return nil
}

type UserUpdatedSuccessPublisher struct{}

func NewFakeUserUpdatedPublisher() publisher.UserUpdatedPublisher {
	return &UserUpdatedSuccessPublisher{}
}

func (p *UserUpdatedSuccessPublisher) Publish(original, new user.User) error {
	return nil
}
