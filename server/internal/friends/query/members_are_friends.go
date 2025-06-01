package query

import (
	"beerbux/internal/friends/db"
	"context"
	"github.com/google/uuid"
)

type MembersAreFriendsQuery struct {
	queries *db.Queries
}

func NewMembersAreFriendsQuery(queries *db.Queries) *MembersAreFriendsQuery {
	return &MembersAreFriendsQuery{
		queries: queries,
	}
}

func (q *MembersAreFriendsQuery) Execute(ctx context.Context, memberID, otherMemberID uuid.UUID) (bool, error) {
	return q.queries.MembersAreFriends(ctx, db.MembersAreFriendsParams{
		MemberID:      memberID,
		OtherMemberID: otherMemberID,
	})
}
