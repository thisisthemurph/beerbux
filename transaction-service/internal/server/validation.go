package server

import (
	"fmt"
	"github.com/google/uuid"
)

func validateStringUUID(s string) error {
	if s == "" {
		return fmt.Errorf("ID is required")
	}
	if _, err := uuid.Parse(s); err != nil {
		return fmt.Errorf("invalid UUID")
	}
	return nil
}
