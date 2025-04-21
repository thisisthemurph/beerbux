package session

import (
	"context"
	"errors"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
)

func memberIsAdmin(ctx context.Context, sessionClient sessionpb.SessionClient, sessionID, memberID string) (bool, error) {
	ssn, err := sessionClient.GetSession(ctx, &sessionpb.GetSessionRequest{
		SessionId: sessionID,
	})

	if err != nil {
		return false, err
	}

	for _, member := range ssn.Members {
		if member.UserId == memberID {
			return member.IsAdmin, nil
		}
	}

	return false, errors.New("member not found")
}
