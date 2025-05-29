package command

import (
	"beerbux/internal/auth/db"
	"beerbux/pkg/otp"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const PasswordResetOTPLength = 6

type InitializePasswordResetCommand struct {
	queries *db.Queries
}

func NewInitializePasswordResetCommand(queries *db.Queries) *InitializePasswordResetCommand {
	return &InitializePasswordResetCommand{
		queries: queries,
	}
}

type InitializePasswordResetResponse struct {
	OTP string
}

func (c *InitializePasswordResetCommand) Execute(ctx context.Context, userID uuid.UUID, password string) (*InitializePasswordResetResponse, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	OTP, err := otp.Generate(PasswordResetOTPLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate one-time password: %w", err)
	}

	err = c.queries.InitializePasswordReset(ctx, db.InitializePasswordResetParams{
		ID: userID,
		PasswordUpdateOtp: sql.NullString{
			String: OTP,
			Valid:  true,
		},
		UpdateHashedPassword: sql.NullString{
			String: string(hashedBytes),
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
