package service

import (
	"beerbux/internal/repository/auth"
	"beerbux/internal/repository/user"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrUserNotFound                 = errors.New("user not found")
	ErrUsernameTaken                = errors.New("username is already taken")
	ErrRefreshTokenNotFound         = errors.New("refresh token not found")
	ErrPasswordMismatch             = errors.New("passwords do not match")
	ErrFailedToGenerateJWT          = errors.New("failed to generate JWT")
	ErrFailedToGenerateRefreshToken = errors.New("failed to generate refresh token")
)

type AuthService struct {
	authRepository *auth.Queries
	userRepository *user.Queries
	options        AuthServiceOptions
}

type AuthServiceOptions struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewAuthService(authRepository *auth.Queries, options AuthServiceOptions) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		options:        options,
	}
}

type LoggedInUserDetails struct {
	ID       uuid.UUID
	Username string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	User         LoggedInUserDetails
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	usr, err := s.authRepository.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.HashedPassword), []byte(password)); err != nil {
		return nil, ErrUserNotFound
	}

	accessToken, err := s.generateJWT(usr.ID, usr.Username, s.options.JWTSecret, s.options.AccessTokenTTL)
	if err != nil {
		return nil, ErrFailedToGenerateJWT
	}

	refreshToken := s.generateRefreshToken()
	hashedRefreshToken, err := s.hashRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrFailedToGenerateRefreshToken
	}

	err = s.authRepository.RegisterRefreshToken(ctx, auth.RegisterRefreshTokenParams{
		UserID:      usr.ID,
		HashedToken: hashedRefreshToken,
		ExpiresAt:   time.Now().Add(s.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, ErrFailedToGenerateRefreshToken
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: LoggedInUserDetails{
			ID:       usr.ID,
			Username: usr.Username,
		},
	}, nil
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Signup(ctx context.Context, name, username, password, verificationPassword string) (*TokenResponse, error) {
	usernameTaken, err := s.userRepository.UserWithUsernameExists(ctx, username)
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

	usr, err := s.authRepository.CreateUser(ctx, auth.CreateUserParams{
		Name:           name,
		Username:       username,
		HashedPassword: string(hashedBytes),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	accessToken, err := s.generateJWT(usr.ID, usr.Username, s.options.JWTSecret, s.options.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	refreshToken := s.generateRefreshToken()
	hashedRefreshToken, err := s.hashRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.authRepository.RegisterRefreshToken(ctx, auth.RegisterRefreshTokenParams{
		UserID:      usr.ID,
		HashedToken: hashedRefreshToken,
		ExpiresAt:   time.Now().Add(s.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateJWT(userID uuid.UUID, username, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"exp":      time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *AuthService) InvalidateRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	userRefreshTokens, err := s.authRepository.GetRefreshTokensByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	for _, t := range userRefreshTokens {
		if err := bcrypt.CompareHashAndPassword([]byte(t.HashedToken), []byte(token)); err == nil {
			if err := s.authRepository.InvalidateRefreshToken(ctx, t.ID); err != nil {
				return fmt.Errorf("failed to invalidate refresh token: %w", err)
			}
			return nil
		}
	}

	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, userID uuid.UUID, token string) (*TokenResponse, error) {
	userRefreshTokens, err := s.authRepository.GetRefreshTokensByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh tokens: %w", err)
	}

	var matchedToken auth.RefreshToken
	var tokenFound bool
	for _, t := range userRefreshTokens {
		if err := bcrypt.CompareHashAndPassword([]byte(t.HashedToken), []byte(token)); err == nil {
			matchedToken = t
			tokenFound = true
			break
		}
	}

	if !tokenFound {
		return nil, ErrRefreshTokenNotFound
	}

	usr, err := s.authRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	newAccessToken, err := s.generateJWT(usr.ID, usr.Username, s.options.JWTSecret, s.options.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	newRefreshToken := s.generateRefreshToken()
	newHashedRefreshToken, err := s.hashRefreshToken(newRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.authRepository.RegisterRefreshToken(ctx, auth.RegisterRefreshTokenParams{
		UserID:      usr.ID,
		HashedToken: newHashedRefreshToken,
		ExpiresAt:   time.Now().Add(s.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	_ = s.authRepository.DeleteRefreshToken(ctx, matchedToken.ID)

	return &TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) generateRefreshToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *AuthService) hashRefreshToken(token string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedToken), nil
}
