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

type UserUpdateResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (c *UpdateUserCommand) Execute(ctx context.Context, userID uuid.UUID, newName string, newUsername string) (*UserUpdateResponse, error) {
	existingUserID, err := c.Queries.GetUserIDByUsername(ctx, newUsername)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to determine if username %s exists", newUsername)
		}
	}
	if err == nil && existingUserID != userID {
		return nil, ErrUsernameExists
	}

	result, err := c.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       userID,
		Name:     newName,
		Username: newUsername,
	})
	if err != nil {
		return nil, err
	}

	return &UserUpdateResponse{
		Name:     result.Name,
		Username: result.Username,
	}, nil
}
