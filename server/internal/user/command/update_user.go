package command

import (
	"beerbux/internal/user/db"
	"context"
	"github.com/google/uuid"
)

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
