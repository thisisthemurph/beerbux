package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/auth-service/internal/producer"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/auth"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/token"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrUsernameExists      = errors.New("username already exists")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
)

type AuthServer struct {
	authpb.UnimplementedAuthServer
	authRepository         *auth.Queries
	authTokenRepository    *token.Queries
	userRegisteredProducer producer.UserRegisteredProducer
	options                AuthServerOptions
	logger                 *slog.Logger
}

type AuthServerOptions struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewAuthServer(
	logger *slog.Logger,
	authRepository *auth.Queries,
	tokenRepository *token.Queries,
	userRegisteredProducer producer.UserRegisteredProducer,
	options AuthServerOptions,
) *AuthServer {
	return &AuthServer{
		authRepository:         authRepository,
		authTokenRepository:    tokenRepository,
		userRegisteredProducer: userRegisteredProducer,
		options:                options,
	}
}

func (s *AuthServer) Login(ctx context.Context, r *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	user, err := s.authRepository.GetUserByUsername(ctx, r.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(r.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := generateJWT(s.options.JWTSecret, user, s.options.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken := generateRefreshToken()
	hashedRefreshToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hashed refresh token: %w", err)
	}

	err = s.authTokenRepository.RegisterRefreshToken(ctx, token.RegisterRefreshTokenParams{
		UserID:      user.ID,
		HashedToken: hashedRefreshToken,
		ExpiresAt:   time.Now().Add(s.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &authpb.UserResponse{
			Id:       user.ID,
			Username: user.Username,
		},
	}, nil
}

func (s *AuthServer) Signup(ctx context.Context, r *authpb.SignupRequest) (*authpb.SignupResponse, error) {
	usernameTaken, err := s.authRepository.UserWithUsernameExists(ctx, r.Username)
	if err != nil {
		return nil, err
	} else if usernameTaken == 1 {
		return nil, ErrUsernameExists
	}

	if r.Password != r.VerificationPassword {
		return nil, ErrPasswordsDoNotMatch
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.authRepository.RegisterUser(ctx, auth.RegisterUserParams{
		ID:             uuid.NewString(),
		Username:       r.Username,
		HashedPassword: string(hashedBytes),
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := generateJWT(s.options.JWTSecret, user, s.options.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken := generateRefreshToken()
	hashedRefreshToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hashed refresh token: %w", err)
	}

	err = s.authTokenRepository.RegisterRefreshToken(ctx, token.RegisterRefreshTokenParams{
		UserID:      user.ID,
		HashedToken: hashedRefreshToken,
		ExpiresAt:   time.Now().Add(s.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, err
	}

	err = s.userRegisteredProducer.Publish(ctx, producer.UserRegisteredEvent{
		UserID:   user.ID,
		Name:     r.Name,
		Username: user.Username,
	})
	if err != nil {
		s.logger.Error("failed to send user registered event", "error", err)
	}

	return &authpb.SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, r *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	userRefreshTokens, err := s.authTokenRepository.GetRefreshTokensByUserID(ctx, r.UserId)
	if err != nil {
		return nil, err
	}

	var matchedToken token.RefreshToken
	var tokenFound bool
	for _, t := range userRefreshTokens {
		if err := bcrypt.CompareHashAndPassword([]byte(t.HashedToken), []byte(r.RefreshToken)); err == nil {
			matchedToken = t
			tokenFound = true
			break
		}
	}

	if !tokenFound {
		return nil, ErrInvalidCredentials
	}

	user, err := s.authRepository.GetUserByID(ctx, r.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	newAccessToken, err := generateJWT(s.options.JWTSecret, user, s.options.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	newRefreshToken := generateRefreshToken()
	newHashedRefreshToken, err := hashRefreshToken(newRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hashed refresh token: %w", err)
	}

	err = s.authTokenRepository.RegisterRefreshToken(ctx, token.RegisterRefreshTokenParams{
		UserID:      user.ID,
		HashedToken: newHashedRefreshToken,
		ExpiresAt:   time.Now().Add(s.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, err
	}

	if err := s.authTokenRepository.DeleteRefreshToken(ctx, matchedToken.ID); err != nil {
		return nil, err
	}

	return &authpb.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
