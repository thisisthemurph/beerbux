package publisher

import (
	"encoding/json"
	"fmt"
	
	"github.com/nats-io/nats.go"
)

const SubjectSessionMemberAdded = "session.member.added"

type SessionMemberAddedEventData struct {
	SessionID string `json:"session_id"`
	MemberID  string `json:"member_id"`
}

type SessionMemberAddedEvent struct {
	Metadata EventMetadata               `json:"metadata"`
	Data     SessionMemberAddedEventData `json:"data"`
}

type SessionMemberAddedPublisher interface {
	Publish(sessionID, memberID string) error
}

type SessionMemberAddedNatsPublisher struct {
	nc      *nats.Conn
	subject string
}

func NewSessionMemberAddedNatsPublisher(nc *nats.Conn) SessionMemberAddedPublisher {
	return &SessionMemberAddedNatsPublisher{
		nc:      nc,
		subject: SubjectSessionMemberAdded,
	}
}

func (p *SessionMemberAddedNatsPublisher) Publish(sessionID, memberID string) error {
	msg := SessionMemberAddedEvent{
		Metadata: NewEventMetadata(SubjectSessionMemberAdded, "1.0.0", sessionID),
		Data: SessionMemberAddedEventData{
			SessionID: sessionID,
			MemberID:  memberID,
		},
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal session member added %v: %w", sessionID, err)
	}

	if err := p.nc.Publish(p.subject, msgData); err != nil {
		return fmt.Errorf("failed to publish %q message: %w", p.subject, err)
	}

	return nil
}
