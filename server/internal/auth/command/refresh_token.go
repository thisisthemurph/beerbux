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
	Queries *db.Queries
	Options config.AuthOptions
}

func NewRefreshTokenCommand(queries *db.Queries, options config.AuthOptions) *RefreshTokenCommand {
	return &RefreshTokenCommand{
		Queries: queries,
		Options: options,
	}
}

func (c *RefreshTokenCommand) Execute(ctx context.Context, userID uuid.UUID, refreshToken string) (*TokenResponse, error) {
	userRefreshTokens, err := c.Queries.GetRefreshTokensByUserID(ctx, userID)
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

	usr, err := c.Queries.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	newAccessToken, err := shared.GenerateJWT(usr.ID, usr.Username, c.Options.JWTSecret, c.Options.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	newRefreshToken := shared.GenerateRefreshToken()
	newHashedRefreshToken, err := shared.HashRefreshToken(newRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = c.Queries.RegisterRefreshToken(ctx, db.RegisterRefreshTokenParams{
		UserID:      usr.ID,
		HashedToken: newHashedRefreshToken,
		ExpiresAt:   time.Now().Add(c.Options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	_ = c.Queries.DeleteRefreshToken(ctx, matchedToken.ID)

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
