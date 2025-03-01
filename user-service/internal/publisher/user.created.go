package publisher

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/pkg/nullish"
)

const SubjectUserCreated = "user.created"

type UserCreatedEventData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

type UserCreatedEvent struct {
	Metadata EventMetadata        `json:"metadata"`
	Data     UserCreatedEventData `json:"user"`
}

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
	msg := UserCreatedEvent{
		Metadata: NewEventMetadata(SubjectUserCreated, "1.0.0", u.ID),
		Data: UserCreatedEventData{
			ID:       u.ID,
			Name:     u.Name,
			Username: u.Username,
			Bio:      nullish.StringOrEmpty(u.Bio),
		},
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal user %v: %w", u.ID, err)
	}

	if err := p.nc.Publish(p.subject, msgData); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.subject, err)
	}

	return nil
}
