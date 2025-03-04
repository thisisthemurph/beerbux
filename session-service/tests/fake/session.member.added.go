package fake

import "github.com/thisisthemurph/beerbux/session-service/internal/publisher"

type sessionMemberAddedPublisher struct{}

func NewFakeSessionMemberAddedPublisher() publisher.SessionMemberAddedPublisher {
	return &sessionMemberAddedPublisher{}
}

func (p *sessionMemberAddedPublisher) Publish(sessionID, memberID string) error {
	return nil
}
