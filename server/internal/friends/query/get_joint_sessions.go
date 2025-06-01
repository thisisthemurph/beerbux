package query

import (
	"beerbux/internal/friends/db"
	"context"
	"github.com/google/uuid"
)

type GetJointSessionsQuery struct {
	queries *db.Queries
}

func NewGetJointSessionsQuery(queries *db.Queries) *GetJointSessionsQuery {
	return &GetJointSessionsQuery{
		queries: queries,
	}
}

func (q *GetJointSessionsQuery) Execute(ctx context.Context, memberID, otherMemberID uuid.UUID) ([]db.Session, error) {
	return q.queries.GetJointSessions(ctx, db.GetJointSessionsParams{
		MemberID:      memberID,
		OtherMemberID: otherMemberID,
	})
}
