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

type GetUserByIdQuery struct {
	Queries *db.Queries
}

func NewGetUserByIdQuery(queries *db.Queries) *GetUserByIdQuery {
	return &GetUserByIdQuery{
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

func (q *GetUserByIdQuery) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
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
