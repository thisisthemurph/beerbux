package publisher

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
)

const SubjectUserCreated = "user.created"

type UserCreatedPublisher interface {
	Publish(u user.User) error
}

type UserCreatedNatsPublisher struct {
	nc      *nats.Conn
	subject string
}

func NewUserCreatedNatsPublisher(nc *nats.Conn) UserCreatedPublisher {
	return &UserCreatedNatsPublisher{
		nc:      nc,
		subject: SubjectUserCreated,
	}
}

func (p *UserCreatedNatsPublisher) Publish(u user.User) error {
	msg, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("failed to marshal user %v: %w", u.ID, err)
	}

	if err := p.nc.Publish(p.subject, msg); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.subject, err)
	}

	return nil
}
