package command

import (
	"beerbux/internal/auth/db"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordResetNotInitialized = errors.New("password reset not initialized")

type ResetPasswordCommand struct {
	queries *db.Queries
}

func NewResetPasswordCommand(queries *db.Queries) *ResetPasswordCommand {
	return &ResetPasswordCommand{
		queries: queries,
	}
}

func (c *ResetPasswordCommand) Execute(ctx context.Context, userEmail, otp, newPassword string) error {
	user, err := c.queries.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return err
	}

	if !user.PasswordUpdateOtp.Valid || !user.PasswordUpdateRequestedAt.Valid {
		return ErrPasswordResetNotInitialized
	}
	if user.UpdateHashedPassword.Valid {
		// This indicates that a password update was initialized, not a password reset
		return ErrPasswordResetNotInitialized
	}

	if user.PasswordUpdateOtp.String != otp {
		return ErrIncorrectOTP
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password: %w", err)
	}

	if err := c.queries.ResetPassword(ctx, db.ResetPasswordParams{
		ID:             user.ID,
		HashedPassword: string(hashedBytes),
	}); err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}
