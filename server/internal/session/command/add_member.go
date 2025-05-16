package command

import (
	"beerbux/internal/common/history"
	"beerbux/internal/session/db"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/pkg/dbtx"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type AddSessionMemberCommand struct {
	dbtx.TX
	Queries              *db.Queries
	SessionHistoryWriter history.SessionHistoryWriter
}

func (cmd *AddSessionMemberCommand) Execute(ctx context.Context, sessionID, memberID, performedByUserID uuid.UUID) error {
	session, err := cmd.Queries.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sessionErr.ErrSessionNotFound
		}
		return fmt.Errorf("failed to fetch session with id %s: %w", sessionID, err)
	}

	if !session.IsActive {
		return sessionErr.ErrCannotUpdateInactiveSession
	}

	err = cmd.Queries.AddMemberToSession(ctx, db.AddMemberToSessionParams{
		SessionID: sessionID,
		MemberID:  memberID,
	})
	if err != nil {
		return fmt.Errorf("failed to add member %s to session %s: %w", memberID, sessionID, err)
	}

	_ = cmd.SessionHistoryWriter.CreateMemberAddedEvent(ctx, sessionID, memberID, performedByUserID)
	return nil
}
