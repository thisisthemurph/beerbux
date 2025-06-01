package query

import (
	"beerbux/internal/friends/db"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/thisisthemurph/fn"
)

type GetFriendsQuery struct {
	queries *db.Queries
}

func NewGetFriendsQuery(queries *db.Queries) *GetFriendsQuery {
	return &GetFriendsQuery{
		queries: queries,
	}
}

type Friend struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Username           string    `json:"username"`
	SharedSessionCount int64     `json:"sharedSessionCount"`
}

func (q *GetFriendsQuery) Execute(ctx context.Context, userID uuid.UUID) ([]Friend, error) {
	friends, err := q.queries.GetFriends(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting friends: %w", err)
	}

	return fn.Map(friends, func(f db.GetFriendsRow) Friend {
		return Friend{
			ID:                 f.ID,
			Name:               f.Name,
			Username:           f.Username,
			SharedSessionCount: f.SharedSessionCount,
		}
	}), nil
}
