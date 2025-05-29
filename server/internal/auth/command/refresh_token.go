package command

import (
	"beerbux/internal/api/config"
	"beerbux/internal/auth/db"
	"beerbux/internal/auth/shared"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type RefreshTokenCommand struct {
	queries *db.Queries
	options config.AuthOptions
}

func NewRefreshTokenCommand(queries *db.Queries, options config.AuthOptions) *RefreshTokenCommand {
	return &RefreshTokenCommand{
		queries: queries,
		options: options,
	}
}

func (c *RefreshTokenCommand) Execute(ctx context.Context, userID uuid.UUID, refreshToken string) (*TokenResponse, error) {
	userRefreshTokens, err := c.queries.GetRefreshTokensByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	var matchedToken db.RefreshToken
	var tokenFound bool
	for _, t := range userRefreshTokens {
		if err := bcrypt.CompareHashAndPassword([]byte(t.HashedToken), []byte(refreshToken)); err == nil {
			matchedToken = t
			tokenFound = true
			break
		}
	}

	if !tokenFound {
		return nil, ErrRefreshTokenNotFound
	}

	usr, err := c.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	newAccessToken, err := shared.GenerateJWT(usr.ID, usr.Username, c.options.JWTSecret, c.options.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	newRefreshToken := shared.GenerateRefreshToken()
	newHashedRefreshToken, err := shared.HashRefreshToken(newRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = c.queries.RegisterRefreshToken(ctx, db.RegisterRefreshTokenParams{
		UserID:      usr.ID,
		HashedToken: newHashedRefreshToken,
		ExpiresAt:   time.Now().Add(c.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	_ = c.queries.DeleteRefreshToken(ctx, matchedToken.ID)

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
