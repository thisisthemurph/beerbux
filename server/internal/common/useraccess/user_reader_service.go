package useraccess

import (
	"beerbux/internal/common/useraccess/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type UserReaderService struct {
	Queries *db.Queries
}

func NewUserReaderService(queries *db.Queries) *UserReaderService {
	return &UserReaderService{
		Queries: queries,
	}
}

type UserAccount struct {
	Debit  float64 `json:"debit"`
	Credit float64 `json:"credit"`
}

type UserResponse struct {
	ID        uuid.UUID   `json:"id"`
	Username  string      `json:"username"`
	Name      string      `json:"name"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Account   UserAccount `json:"account"`
}

func (q *UserReaderService) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	usr, err := q.Queries.GetUserByID(ctx, userID)
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
		Account: UserAccount{
			Debit:  usr.Debit,
			Credit: usr.Credit,
		},
	}, nil
}

func (q *UserReaderService) GetUserByUsername(ctx context.Context, username string) (*UserResponse, error) {
	usr, err := q.Queries.GetByUsername(ctx, username)
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
		Account: UserAccount{
			Debit:  usr.Debit,
			Credit: usr.Credit,
		},
	}, nil
}

func (q *UserReaderService) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {
	usr, err := q.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user with email %s: %w", email, err)
	}

	return &UserResponse{
		ID:        usr.ID,
		Username:  usr.Username,
		Name:      usr.Name,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Account: UserAccount{
			Debit:  usr.Debit,
			Credit: usr.Credit,
		},
	}, nil
}

func (q *UserReaderService) UserWithUsernameExists(ctx context.Context, username string) (bool, error) {
	return q.Queries.UserWithUsernameExists(ctx, username)
}

func (q *UserReaderService) UserWithEmailExists(ctx context.Context, username string) (bool, error) {
	return q.Queries.UserWithEmailExists(ctx, username)
}
