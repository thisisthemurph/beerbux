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
	"regexp"
	"strings"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type LoginCommand struct {
	queries *db.Queries
	options config.AuthOptions
}

func NewLoginCommand(queries *db.Queries, options config.AuthOptions) *LoginCommand {
	return &LoginCommand{
		queries: queries,
		options: options,
	}
}

type LoggedInUserDetails struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
}

type LoginResponse struct {
	AccessToken  string              `json:"accessToken"`
	RefreshToken string              `json:"refreshToken"`
	User         LoggedInUserDetails `json:"user"`
}

func (c *LoginCommand) Execute(ctx context.Context, usernameOrEmail, password string) (*LoginResponse, error) {
	user, err := c.getUserByUsernameOrEmail(ctx, usernameOrEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, ErrUserNotFound
	}

	accessToken, err := shared.GenerateJWT(user.ID, user.Username, user.Email, c.options.JWTSecret, c.options.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	refreshToken := shared.GenerateRefreshToken()
	hashedRefreshToken, err := shared.HashRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = c.queries.RegisterRefreshToken(ctx, db.RegisterRefreshTokenParams{
		UserID:      user.ID,
		HashedToken: hashedRefreshToken,
		ExpiresAt:   time.Now().Add(c.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: LoggedInUserDetails{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Name:     user.Name,
		},
	}, nil
}

func (c *LoginCommand) getUserByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (db.User, error) {
	if isEmail(usernameOrEmail) {
		return c.queries.GetUserByEmail(ctx, usernameOrEmail)
	}
	return c.queries.GetUserByUsername(ctx, usernameOrEmail)
}

var (
	userRegexp    = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp    = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
	userDotRegexp = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
)

func isEmail(email string) bool {
	if len(email) < 6 || len(email) > 254 {
		return false
	}
	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return false
	}
	user := email[:at]
	host := email[at+1:]
	if len(user) > 64 {
		return false
	}
	if userDotRegexp.MatchString(user) || !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return false
	}

	return true
}
