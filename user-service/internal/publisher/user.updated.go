package publisher

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
)

const SubjectUserUpdated = "user.updated"

type UserUpdatedPublisher interface {
	Publish(original user.User, new user.User) error
}

type UserUpdatedNatsPublisher struct {
	nc      *nats.Conn
	subject string
}

func NewUserUpdatedNatsPublisher(nc *nats.Conn) UserUpdatedPublisher {
	return &UserUpdatedNatsPublisher{
		nc:      nc,
		subject: SubjectUserUpdated,
	}
}

func (p *UserUpdatedNatsPublisher) Publish(original user.User, new user.User) error {
	msg := struct {
		Original user.User `json:"original"`
		New      user.User `json:"new"`
	}{
		Original: original,
		New:      new,
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := p.nc.Publish(p.subject, msgData); err != nil {
		return err
	}

	return nil
}
