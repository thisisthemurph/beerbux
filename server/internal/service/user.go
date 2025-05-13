package service

import (
	"beerbux/internal/repository/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	userRepository *user.Queries
}

func NewUserService(userRepository *user.Queries) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

type UserResponse struct {
	ID        uuid.UUID
	Username  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Account   struct {
		Debit  float64
		Credit float64
	}
}

func (s *UserService) Get(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	usr, err := s.userRepository.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user with id %s: %w", userID, err)
	}

	return &UserResponse{
		ID:        usr.ID,
		Username:  usr.Username,
		Name:      usr.Name,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Account: struct {
			Debit  float64
			Credit float64
		}{
			Debit:  usr.Debit,
			Credit: usr.Credit,
		},
	}, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*UserResponse, error) {
	usr, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user with username %s: %w", username, err)
	}

	return &UserResponse{
		ID:        usr.ID,
		Username:  usr.Username,
		Name:      usr.Name,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Account: struct {
			Debit  float64
			Credit float64
		}{
			Debit:  usr.Debit,
			Credit: usr.Credit,
		},
	}, nil
}
