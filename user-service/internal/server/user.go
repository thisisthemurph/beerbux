package server

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
)

type UserServer struct {
	userpb.UnimplementedUserServer
	userRepository *user.Queries
}

func NewUserServer(userRepository *user.Queries) *UserServer {
	return &UserServer{
		userRepository: userRepository,
	}
}

func (s *UserServer) GetUser(ctx context.Context, r *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	if err := validateGetUserRequest(r); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	u, err := s.userRepository.GetUser(ctx, r.UserId)
	if err != nil {
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
		Bio:      createNullString(r.Bio),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user %v: %w", r.Username, err)
	}

	return &userpb.UserResponse{
		UserId:   u.ID,
		Username: u.Username,
		Bio:      parseNullString(u.Bio),
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, r *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	if err := validateUpdateUserRequest(r); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	u, err := s.userRepository.UpdateUser(ctx, user.UpdateUserParams{
		ID:       r.UserId,
		Username: r.Username,
		Bio:      createNullString(r.Bio),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update user %v: %w", r.UserId, err)
	}

	return &userpb.UserResponse{
		UserId:   u.ID,
		Username: u.Username,
		Bio:      parseNullString(u.Bio),
	}, nil
}
