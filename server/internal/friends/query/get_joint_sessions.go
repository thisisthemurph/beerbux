package query

import (
	"beerbux/internal/friends/db"
	"context"
	"github.com/google/uuid"
)

type GetJointSessionIDsQuery struct {
	queries *db.Queries
}

func NewGetJointSessionIDsQuery(queries *db.Queries) *GetJointSessionIDsQuery {
	return &GetJointSessionIDsQuery{
		queries: queries,
	}
}

func (q *GetJointSessionIDsQuery) Execute(ctx context.Context, memberID, otherMemberID uuid.UUID) ([]uuid.UUID, error) {
	return q.queries.GetJointSessionIDs(ctx, db.GetJointSessionIDsParams{
		MemberID:      memberID,
		OtherMemberID: otherMemberID,
	})
}
