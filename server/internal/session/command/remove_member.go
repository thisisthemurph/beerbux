package command

import (
	"beerbux/internal/session/db"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/history"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type RemoveSessionMemberCommand struct {
	Queries              *db.Queries
	SessionHistoryWriter history.SessionHistoryWriter
}

func NewRemoveSessionMemberCommand(queries *db.Queries, historyWriter history.SessionHistoryWriter) *RemoveSessionMemberCommand {
	return &RemoveSessionMemberCommand{
		Queries:              queries,
		SessionHistoryWriter: historyWriter,
	}
}

func (cmd *RemoveSessionMemberCommand) Execute(ctx context.Context, sessionID, memberID uuid.UUID, performedByUserID uuid.UUID) error {
	if exists, err := cmd.Queries.SessionExists(ctx, sessionID); err != nil {
		return fmt.Errorf("could not check if session exists: %w", err)
	} else if !exists {
		return sessionErr.ErrSessionNotFound
	}

	memberToRemove, err := cmd.Queries.GetSessionMember(ctx, db.GetSessionMemberParams{
		SessionID: sessionID,
		ID:        memberID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sessionErr.ErrSessionMemberNotFound
		}
		return fmt.Errorf("could not get session member: %w", err)
	}

	if count, err := cmd.Queries.CountSessionMembers(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to count session members: %w", err)
	} else if count == 1 {
		return sessionErr.ErrSessionMustHaveAtLeastOneMember
	}

	if memberToRemove.IsAdmin {
		if count, err := cmd.Queries.CountSessionAdminMembers(ctx, sessionID); err != nil {
			return fmt.Errorf("could not get session admin member count: %w", err)
		} else if count <= 1 {
			return sessionErr.ErrSessionMustHaveAtLeastOneAdmin
		}
	}

	if err := cmd.Queries.DeleteSessionMember(ctx, db.DeleteSessionMemberParams{
		SessionID: sessionID,
		MemberID:  memberToRemove.ID,
	}); err != nil {
		return fmt.Errorf("failed to remove member from the session: %w", err)
	}

	if performedByUserID == memberToRemove.ID {
		_ = cmd.SessionHistoryWriter.CreateMemberLeftEvent(ctx, sessionID, memberToRemove.ID)
	} else {
		_ = cmd.SessionHistoryWriter.CreateMemberRemovedEvent(ctx, sessionID, memberToRemove.ID, performedByUserID)
	}

	return nil
}
