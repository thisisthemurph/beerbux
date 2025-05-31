package command

import (
	"beerbux/internal/auth/db"
	"beerbux/pkg/otp"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

const PasswordResetOTPLength int = 6

type InitializePasswordResetCommand struct {
	queries *db.Queries
}

type InitializePasswordResetResponse struct {
	OTP string
}

func NewInitializePasswordResetCommand(queries *db.Queries) *InitializePasswordResetCommand {
	return &InitializePasswordResetCommand{
		queries: queries,
	}
}

func (c *InitializePasswordResetCommand) Execute(ctx context.Context, userID uuid.UUID) (*InitializePasswordResetResponse, error) {
	OTP, err := otp.Generate(PasswordResetOTPLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}

	err = c.queries.InitializePasswordReset(ctx, db.InitializePasswordResetParams{
		ID: userID,
		PasswordUpdateOtp: sql.NullString{
			String: OTP,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize password reset: %w", err)
	}

	return &InitializePasswordResetResponse{
		OTP: OTP,
	}, nil
}
