package publisher

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
)

const SubjectUserUpdated = "user.updated"

type UserUpdatedEventData struct {
	UserID        string                 `json:"user_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"`
}

type UserUpdatedEvent struct {
	Metadata EventMetadata        `json:"metadata"`
	Data     UserUpdatedEventData `json:"user"`
}

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
	msg := UserUpdatedEvent{
		Metadata: NewEventMetadata(SubjectUserUpdated, "1.0.0", new.ID),
		Data: UserUpdatedEventData{
			UserID:        new.ID,
			UpdatedFields: p.determineUpdatedFields(original, new),
		},
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal user %v: %w", new.ID, err)
	}

	if err := p.nc.Publish(p.subject, msgData); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.subject, err)
	}

	return nil
}

func (p *UserUpdatedNatsPublisher) determineUpdatedFields(original, new user.User) map[string]interface{} {
	updatedFields := make(map[string]interface{})

	if original.Username != new.Username {
		updatedFields["username"] = new.Username
	}

	if original.Bio != new.Bio {
		updatedFields["bio"] = new.Bio
	}

	updatedFields["updated_at"] = new.UpdatedAt

	return updatedFields
}
