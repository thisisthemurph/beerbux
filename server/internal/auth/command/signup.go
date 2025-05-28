package command

import (
	"beerbux/internal/api/config"
	"beerbux/internal/auth/db"
	"beerbux/internal/auth/shared"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrUsernameTaken    = errors.New("username already taken")
	ErrPasswordMismatch = errors.New("password mismatch")
)

type SignupCommand struct {
	Queries *db.Queries
	Options config.AuthOptions
}

func NewSignupCommand(queries *db.Queries, options config.AuthOptions) *SignupCommand {
	return &SignupCommand{
		Queries: queries,
		Options: options,
	}
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

func (c *SignupCommand) Execute(ctx context.Context, name, username, email, password, verificationPassword string) (*TokenResponse, error) {
	usernameTaken, err := c.Queries.UserWithUsernameExists(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to determine if username exists: %w", err)
	} else if usernameTaken {
		return nil, ErrUsernameTaken
	}

	if password != verificationPassword {
		return nil, ErrPasswordMismatch
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	usr, err := c.Queries.CreateUser(ctx, db.CreateUserParams{
		Name:           name,
		Username:       username,
		Email:          email,
		HashedPassword: string(hashedBytes),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
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

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
