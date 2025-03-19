package server

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/auth-service/internal/producer"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/auth"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/token"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound         = status.Error(codes.NotFound, "user not found")
	ErrUsernameTaken        = status.Error(codes.InvalidArgument, "username is already taken")
	ErrRefreshTokenNotFound = status.Error(codes.Unauthenticated, "refresh token not found")
	ErrPasswordMismatch     = status.Error(codes.InvalidArgument, "passwords do not match")
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
			return nil, ErrUserNotFound
		}
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(r.Password)); err != nil {
		return nil, ErrUserNotFound
	}

	accessToken, err := generateJWT(s.options.JWTSecret, user, s.options.AccessTokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate JWT: %v", err)
	}

	refreshToken := generateRefreshToken()
	hashedRefreshToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate hashed refresh token: %v", err)
	}

	err = s.authTokenRepository.RegisterRefreshToken(ctx, token.RegisterRefreshTokenParams{
		UserID:      user.ID,
		HashedToken: hashedRefreshToken,
		ExpiresAt:   time.Now().Add(s.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store refresh token: %v", err)
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
		return nil, ErrUsernameTaken
	}

	if r.Password != r.VerificationPassword {
		return nil, ErrPasswordMismatch
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	user, err := s.authRepository.RegisterUser(ctx, auth.RegisterUserParams{
		ID:             uuid.NewString(),
		Username:       r.Username,
		HashedPassword: string(hashedBytes),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	accessToken, err := generateJWT(s.options.JWTSecret, user, s.options.AccessTokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate JWT: %v", err)
	}

	refreshToken := generateRefreshToken()
	hashedRefreshToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate hashed refresh token: %v", err)
	}

	err = s.authTokenRepository.RegisterRefreshToken(ctx, token.RegisterRefreshTokenParams{
		UserID:      user.ID,
		HashedToken: hashedRefreshToken,
		ExpiresAt:   time.Now().Add(s.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store refresh token: %v", err)
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
		return nil, ErrRefreshTokenNotFound
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
		return nil, ErrRefreshTokenNotFound
	}

	user, err := s.authRepository.GetUserByID(ctx, r.UserId)
	if err != nil {
		return nil, ErrUserNotFound
	}

	newAccessToken, err := generateJWT(s.options.JWTSecret, user, s.options.AccessTokenTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate JWT: %v", err)
	}

	newRefreshToken := generateRefreshToken()
	newHashedRefreshToken, err := hashRefreshToken(newRefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate hashed refresh token: %v", err)
	}

	err = s.authTokenRepository.RegisterRefreshToken(ctx, token.RegisterRefreshTokenParams{
		UserID:      user.ID,
		HashedToken: newHashedRefreshToken,
		ExpiresAt:   time.Now().Add(s.options.RefreshTokenTTL),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store refresh token: %v", err)
	}

	if err := s.authTokenRepository.DeleteRefreshToken(ctx, matchedToken.ID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete refresh token: %v", err)
	}

	return &authpb.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
