package command

import (
	"beerbux/internal/user/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var ErrUsernameExists = errors.New("username exists")

type UpdateUserCommand struct {
	Queries *db.Queries
}

func NewUpdateUserCommand(queries *db.Queries) *UpdateUserCommand {
	return &UpdateUserCommand{
		Queries: queries,
	}
}

func (c *UpdateUserCommand) Execute(ctx context.Context, userID uuid.UUID, newName string, newUsername string) error {
	existingUserID, err := c.Queries.GetUserIDByUsername(ctx, newUsername)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to determine if username %s exists", newUsername)
		}
	}
	if err == nil && existingUserID != userID {
		return ErrUsernameExists
	}

	if err := c.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       userID,
		Name:     newName,
		Username: newUsername,
	}); err != nil {
		return err
	}

	return nil
}
