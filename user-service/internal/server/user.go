package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"

	"github.com/thisisthemurph/beerbux/user-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/pkg/nullish"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
)

type UserServer struct {
	userpb.UnimplementedUserServer
	userRepository       *user.Queries
	userCreatedPublisher publisher.UserCreatedPublisher
	userUpdatedPublisher publisher.UserUpdatedPublisher
	logger               *slog.Logger
}

func NewUserServer(
	userRepository *user.Queries,
	userCreatedPublisher publisher.UserCreatedPublisher,
	userUpdatedPublisher publisher.UserUpdatedPublisher,
	logger *slog.Logger,
) *UserServer {
	return &UserServer{
		userRepository:       userRepository,
		userCreatedPublisher: userCreatedPublisher,
		userUpdatedPublisher: userUpdatedPublisher,
		logger:               logger,
	}
}

func (s *UserServer) GetUser(ctx context.Context, r *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	if err := validateGetUserRequest(r); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	u, err := s.userRepository.GetUser(ctx, r.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", r.UserId)
		}
		s.logger.Error("failed to get user", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.GetUserResponse{
		UserId:     u.ID,
		Name:       u.Name,
		Username:   u.Username,
		Bio:        nullish.ParseNullString(u.Bio),
		NetBalance: u.Net,
	}, nil
}

func (s *UserServer) GetUserByUsername(ctx context.Context, r *userpb.GetUserByUsernameRequest) (*userpb.GetUserResponse, error) {
	u, err := s.userRepository.GetUserByUsername(ctx, r.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user with username %s not found", r.Username)
		}
		s.logger.Error("failed to get user", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.GetUserResponse{
		UserId:     u.ID,
		Name:       u.Name,
		Username:   u.Username,
		Bio:        nullish.ParseNullString(u.Bio),
		NetBalance: u.Net,
	}, nil
}

func (s *UserServer) GetUserBalance(ctx context.Context, r *userpb.GetUserRequest) (*userpb.UserBalanceResponse, error) {
	balance, err := s.userRepository.GetUserBalances(ctx, r.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", r.UserId)
		}
		s.logger.Error("failed to get user balance", "error", err)
		return nil, status.Error(codes.Internal, "failed to get user balance")
	}

	return &userpb.UserBalanceResponse{
		Credit: balance.Credit,
		Debit:  balance.Debit,
		Net:    balance.Net,
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, r *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	if err := validateUpdateUserRequest(r); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	originalUser, err := s.userRepository.GetUser(ctx, r.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", r.UserId)
		}
		s.logger.Error("failed to get user", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	u, err := s.userRepository.UpdateUser(ctx, user.UpdateUserParams{
		ID:       r.UserId,
		Name:     r.Name,
		Username: r.Username,
		Bio:      nullish.CreateNullString(r.Bio),
	})

	if err != nil {
		s.logger.Error("failed to update user", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := &userpb.UserResponse{
		UserId:   u.ID,
		Name:     u.Name,
		Username: u.Username,
		Bio:      nullish.ParseNullString(u.Bio),
	}

	if err := s.userUpdatedPublisher.Publish(originalUser, u); err != nil {
		s.logger.Error("failed to publish user updated event", "error", err)
		return result, status.Errorf(codes.Internal, "failed to publish user updated event: %v", err)
	}

	return result, nil
}
