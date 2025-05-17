package useraccess

import (
	"beerbux/internal/common/useraccess/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type GetUserByUsernameQuery struct {
	Queries *db.Queries
}

func NewGetUserByUsernameQuery(queries *db.Queries) *GetUserByUsernameQuery {
	return &GetUserByUsernameQuery{
		Queries: queries,
	}
}

func (q *GetUserByUsernameQuery) Execute(ctx context.Context, username string) (*UserResponse, error) {
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
