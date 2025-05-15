package command

import (
	"beerbux/internal/session/db"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/history"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type UpdateSessionActiveStateCommand struct {
	Queries              *db.Queries
	SessionHistoryWriter history.SessionHistoryWriter
}

func NewUpdateSessionActionStateCommand(queries *db.Queries, sessionHistoryWriter history.SessionHistoryWriter) *UpdateSessionActiveStateCommand {
	return &UpdateSessionActiveStateCommand{
		Queries:              queries,
		SessionHistoryWriter: sessionHistoryWriter,
	}
}

func (cmd *UpdateSessionActiveStateCommand) Execute(ctx context.Context, sessionID, performedByUserID uuid.UUID, isActive bool) error {
	exists, err := cmd.Queries.SessionExists(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to determine if session exists: %w", err)
	}
	if !exists {
		return sessionErr.ErrSessionNotFound
	}

	err = cmd.Queries.UpdateSessionActiveState(ctx, db.UpdateSessionActiveStateParams{
		ID:       sessionID,
		IsActive: isActive,
	})
	if err != nil {
		return fmt.Errorf("failed to update session active state: %w", err)
	}

	if isActive {
		_ = cmd.SessionHistoryWriter.CreateSessionOpenedEvent(ctx, sessionID, performedByUserID)
	} else {
		_ = cmd.SessionHistoryWriter.CreateSessionClosedEvent(ctx, sessionID, performedByUserID)
	}

	return nil
}
