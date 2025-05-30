package command

import (
	"beerbux/internal/auth/db"
	"beerbux/pkg/otp"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

const EmailUpdateOTPLength = 6

type InitializeUpdateEmailCommand struct {
	queries *db.Queries
}

func NewInitializeUpdateEmailCommand(queries *db.Queries) *InitializeUpdateEmailCommand {
	return &InitializeUpdateEmailCommand{
		queries: queries,
	}
}

type InitializeUpdateEmailCommandResponse struct {
	OTP string
}

func (c *InitializeUpdateEmailCommand) Execute(ctx context.Context, userID uuid.UUID, email string) (*InitializeUpdateEmailCommandResponse, error) {
	OTP, err := otp.Generate(EmailUpdateOTPLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate one-time password: %w", err)
	}

	err = c.queries.InitialiseUpdateEmail(ctx, db.InitialiseUpdateEmailParams{
		ID: userID,
		UpdateEmail: sql.NullString{
			String: email,
			Valid:  true,
		},
		EmailUpdateOtp: sql.NullString{
			String: OTP,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialise update email: %w", err)
	}

	return &InitializeUpdateEmailCommandResponse{
		OTP: OTP,
	}, nil
}
