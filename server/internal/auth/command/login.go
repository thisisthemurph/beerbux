package command

import (
	"beerbux/internal/api/config"
	"beerbux/internal/auth/db"
	"beerbux/internal/auth/shared"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type LoginCommand struct {
	Queries *db.Queries
	Options config.AuthOptions
}

func NewLoginCommand(queries *db.Queries, options config.AuthOptions) *LoginCommand {
	return &LoginCommand{
		Queries: queries,
		Options: options,
	}
}

type LoggedInUserDetails struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
}

type LoginResponse struct {
	AccessToken  string              `json:"accessToken"`
	RefreshToken string              `json:"refreshToken"`
	User         LoggedInUserDetails `json:"user"`
}

func (c *LoginCommand) Execute(ctx context.Context, username, password string) (*LoginResponse, error) {
	usr, err := c.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.HashedPassword), []byte(password)); err != nil {
		return nil, ErrUserNotFound
	}

	accessToken, err := shared.GenerateJWT(usr.ID, usr.Username, c.Options.JWTSecret, c.Options.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	refreshToken := shared.GenerateRefreshToken()
	hashedRefreshToken, err := shared.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = c.Queries.RegisterRefreshToken(ctx, db.RegisterRefreshTokenParams{
		UserID:      usr.ID,
		HashedToken: hashedRefreshToken,
		ExpiresAt:   time.Now().Add(c.Options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: LoggedInUserDetails{
			ID:       usr.ID,
			Username: usr.Username,
			Name:     usr.Name,
		},
	}, nil
}
