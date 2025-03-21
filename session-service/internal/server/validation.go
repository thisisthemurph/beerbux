package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
)

func validateGetSessionRequest(r *sessionpb.GetSessionRequest) error {
	return validateStringUUID(r.SessionId)
}

func validateCreateSessionRequest(r *sessionpb.CreateSessionRequest) error {
	if err := validateStringUUID(r.UserId); err != nil {
		return err
	}
	return validateSessionName(r.Name)
}

func validateListSessionsForUserRequest(r *sessionpb.ListSessionsForUserRequest) error {
	return validateStringUUID(r.UserId)
}

func validateSessionName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) < 3 {
		return fmt.Errorf("name must be at least 3 characters")
	}
	if len(name) > 20 {
		return fmt.Errorf("name must be at most 20 characters")
	}
	return nil
}

func validateStringUUID(s string) error {
	if s == "" {
		return fmt.Errorf("ID is required")
	}
	if _, err := uuid.Parse(s); err != nil {
		return fmt.Errorf("invalid UUID")
	}
	return nil
}
