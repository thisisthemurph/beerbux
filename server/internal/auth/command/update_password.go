package command

import (
	"beerbux/internal/auth/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const OTPTimeToLiveMinutes int = 30

var ErrPasswordResetNotInitialized = errors.New("password reset not initialized")
var ErrOTPExpired = errors.New("OTP has expired")
var ErrIncorrectOTP = errors.New("incorrect OTP")

type UpdatePasswordCommand struct {
	queries *db.Queries
}

func NewResetPasswordCommand(queries *db.Queries) *UpdatePasswordCommand {
	return &UpdatePasswordCommand{
		queries: queries,
	}
}

func (c *UpdatePasswordCommand) Execute(ctx context.Context, userID uuid.UUID, OTP string) error {
	user, err := c.queries.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user for password reset: %w", err)
	}

	if !user.PasswordUpdateRequestedAt.Valid || !user.PasswordUpdateOtp.Valid {
		return ErrPasswordResetNotInitialized
	}

	now := time.Now()
	expirationTime := user.PasswordUpdateRequestedAt.Time.Add(time.Duration(OTPTimeToLiveMinutes) * time.Minute)
	if expirationTime.Before(now) {
		return ErrOTPExpired
	}

	if OTP != user.PasswordUpdateOtp.String {
		return ErrIncorrectOTP
	}

	if err := c.queries.UpdatePassword(ctx, userID); err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}
