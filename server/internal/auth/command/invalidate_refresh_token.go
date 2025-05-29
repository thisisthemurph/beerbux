package command

import (
	"beerbux/internal/auth/db"
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type InvalidateRefreshTokenCommand struct {
	queries *db.Queries
}

func NewInvalidateRefreshTokenCommand(authRepository *db.Queries) *InvalidateRefreshTokenCommand {
	return &InvalidateRefreshTokenCommand{
		queries: authRepository,
	}
}

func (c *InvalidateRefreshTokenCommand) Execute(ctx context.Context, userID uuid.UUID, token string) error {
	userRefreshTokens, err := c.queries.GetRefreshTokensByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	for _, t := range userRefreshTokens {
		if err := bcrypt.CompareHashAndPassword([]byte(t.HashedToken), []byte(token)); err == nil {
			if err := c.queries.InvalidateRefreshToken(ctx, t.ID); err != nil {
				return fmt.Errorf("failed to invalidate refresh token: %w", err)
			}
			return nil
		}
	}

	return nil
}
