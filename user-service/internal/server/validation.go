package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
)

func validateGetUserRequest(r *userpb.GetUserRequest) error {
	return validateStringUUID(r.UserId)
}

func validateCreateUserRequest(r *userpb.CreateUserRequest) error {
	if err := validateUserName(r.Name); err != nil {
		return err
	}
	return validateUsername(r.Username)
}

func validateUpdateUserRequest(r *userpb.UpdateUserRequest) error {
	if err := validateUserName(r.Name); err != nil {
		return err
	}
	if err := validateStringUUID(r.UserId); err != nil {
		return err
	}
	return validateUsername(r.Username)
}

func validateUserName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 3 characters")
	}
	if len(name) > 20 {
		return fmt.Errorf("name must be at most 20 characters")
	}
	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if len(username) > 20 {
		return fmt.Errorf("username must be at most 20 characters")
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
