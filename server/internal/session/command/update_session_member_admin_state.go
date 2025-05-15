package command

import (
	"beerbux/internal/session/db"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/history"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type UpdateSessionMemberAdminStateCommand struct {
	Queries              *db.Queries
	SessionHistoryWriter history.SessionHistoryWriter
}

func NewUpdateSessionMemberAdminStateCommand(queries *db.Queries, historyWriter history.SessionHistoryWriter) *UpdateSessionMemberAdminStateCommand {
	return &UpdateSessionMemberAdminStateCommand{
		Queries:              queries,
		SessionHistoryWriter: historyWriter,
	}
}

func (cmd *UpdateSessionMemberAdminStateCommand) Execute(
	ctx context.Context,
	sessionID,
	memberID uuid.UUID,
	isAdmin bool,
	performedByUserID uuid.UUID,
) error {
	createHistoryEvent := cmd.SessionHistoryWriter.CreateMemberPromotedToAdminEvent
	if !isAdmin {
		createHistoryEvent = cmd.SessionHistoryWriter.CreateMemberDemotedFromAdminEvent
		if count, err := cmd.Queries.CountSessionAdminMembers(ctx, sessionID); err != nil {
			return err
		} else if count <= 1 {
			return sessionErr.ErrSessionMustHaveAtLeastOneAdmin
		}
	}

	if err := cmd.Queries.UpdateSessionMemberAdminState(ctx, db.UpdateSessionMemberAdminStateParams{
		SessionID: sessionID,
		MemberID:  memberID,
		IsAdmin:   isAdmin,
	}); err != nil {
		return fmt.Errorf("failed to update session member admin state: %w", err)
	}

	_ = createHistoryEvent(ctx, sessionID, memberID, performedByUserID)
	return nil
}
