package handler

import "errors"

// validatePassword validates the given password and returns an error if validation fails.
func validatePassword(p string) error {
	if len(p) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}
