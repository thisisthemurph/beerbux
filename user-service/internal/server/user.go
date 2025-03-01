package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
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

func (s *UserServer) GetUser(ctx context.Context, r *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	if err := validateGetUserRequest(r); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	u, err := s.userRepository.GetUser(ctx, r.UserId)
	if err != nil {
		s.logger.Error("failed to get user", "error", err)
		return nil, fmt.Errorf("failed to get user %v: %w", r.UserId, err)
	}

	return &userpb.UserResponse{
		UserId:   u.ID,
		Username: u.Username,
	}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, r *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	if err := validateCreateUserRequest(r); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	u, err := s.userRepository.CreateUser(ctx, user.CreateUserParams{
		ID:       uuid.New().String(),
		Username: r.Username,
		Bio:      nullish.CreateNullString(r.Bio),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user %v: %w", r.Username, err)
	}

	if err := s.userCreatedPublisher.Publish(u); err != nil {
		s.logger.Error("failed to publish user created event", "error", err)
		return nil, err
	}

	return &userpb.UserResponse{
		UserId:   u.ID,
		Username: u.Username,
		Bio:      nullish.ParseNullString(u.Bio),
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, r *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	if err := validateUpdateUserRequest(r); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	originalUser, err := s.userRepository.GetUser(ctx, r.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %v: %w", r.UserId, err)
	}

	u, err := s.userRepository.UpdateUser(ctx, user.UpdateUserParams{
		ID:       r.UserId,
		Username: r.Username,
		Bio:      nullish.CreateNullString(r.Bio),
	})

	if err != nil {
		s.logger.Error("failed to update user", "error", err)
		return nil, fmt.Errorf("failed to update user %v: %w", r.UserId, err)
	}

	result := &userpb.UserResponse{
		UserId:   u.ID,
		Username: u.Username,
		Bio:      nullish.ParseNullString(u.Bio),
	}

	if err := s.userUpdatedPublisher.Publish(originalUser, u); err != nil {
		s.logger.Error("failed to publish user updated event", "error", err)
		return result, err
	}

	return result, nil
}
