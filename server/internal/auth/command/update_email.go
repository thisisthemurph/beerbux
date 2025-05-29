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

type UpdateEmailCommand struct {
	queries *db.Queries
}

func NewUpdateEmailCommand(queries *db.Queries) *UpdateEmailCommand {
	return &UpdateEmailCommand{
		queries: queries,
	}
}

func (c *UpdateEmailCommand) Execute(ctx context.Context, userID uuid.UUID, OTP string) error {
	user, err := c.queries.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user for email update: %w", err)
	}

	if !user.EmailUpdateRequestedAt.Valid || !user.EmailUpdateOtp.Valid {
		return ErrProcessNotInitialized
	}

	expirationTime := user.EmailUpdateRequestedAt.Time.Add(time.Duration(OTPTimeToLiveMinutes) * time.Minute)
	if expirationTime.Before(time.Now()) {
		return ErrOTPExpired
	}

	if OTP != user.EmailUpdateOtp.String {
		return ErrIncorrectOTP
	}

	if err := c.queries.UpdateEmail(ctx, userID); err != nil {
		return fmt.Errorf("failed to update email: %w", err)
	}
	return nil
}
