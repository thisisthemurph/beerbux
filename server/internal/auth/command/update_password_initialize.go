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

const PasswordUpdateOTPLength = 6

type InitializeUpdatePasswordCommand struct {
	queries *db.Queries
}

func NewInitializeUpdatePasswordCommand(queries *db.Queries) *InitializeUpdatePasswordCommand {
	return &InitializeUpdatePasswordCommand{
		queries: queries,
	}
}

type InitializeUpdatePasswordResponse struct {
	OTP string
}

func (c *InitializeUpdatePasswordCommand) Execute(ctx context.Context, userID uuid.UUID, password string) (*InitializeUpdatePasswordResponse, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	OTP, err := otp.Generate(PasswordUpdateOTPLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate one-time password: %w", err)
	}

	err = c.queries.InitializePasswordUpdate(ctx, db.InitializePasswordUpdateParams{
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
		return nil, fmt.Errorf("failed to initialize password update: %w", err)
	}

	return &InitializeUpdatePasswordResponse{
		OTP: OTP,
	}, nil
}
