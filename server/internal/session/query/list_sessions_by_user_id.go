package query

import (
	"beerbux/internal/session/db"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type ListSessionsByUserIDQuery struct {
	Queries *db.Queries
}

func NewListSessionsByUserIDQuery(queries *db.Queries) *ListSessionsByUserIDQuery {
	return &ListSessionsByUserIDQuery{
		Queries: queries,
	}
}

func (q *ListSessionsByUserIDQuery) Execute(ctx context.Context, userID uuid.UUID, limit int32) ([]SessionResponse, error) {
	rows, err := q.Queries.ListSessionsForUser(ctx, db.ListSessionsForUserParams{
		MemberID: userID,
		Column2:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions for user %s: %w", userID, err)
	}

	sessions := make([]SessionResponse, 0, len(rows)/2)
	// Use a map to check of we have seen a session row already
	sessionsMap := make(map[uuid.UUID]SessionResponse, len(rows)/2)

	for _, r := range rows {
		s, exists := sessionsMap[r.ID]
		if !exists {
			s = SessionResponse{
				ID:       r.ID,
				Name:     r.Name,
				IsActive: r.IsActive,
				Members:  make([]SessionMember, 0, 4),
				Total:    r.TotalAmount,
			}
			sessionsMap[r.ID] = s
			sessions = append(sessions, s)
		}

		s.Members = append(s.Members, SessionMember{
			ID:       r.MemberID,
			Name:     r.MemberName,
			Username: r.MemberUsername,
		})
	}

	return sessions, nil
}
