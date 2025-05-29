package command

import (
	"beerbux/internal/auth/db"
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type ComparePasswordCommand struct {
	queries *db.Queries
}

func NewComparePasswordCommand(queries *db.Queries) *ComparePasswordCommand {
	return &ComparePasswordCommand{
		queries: queries,
	}
}

func (c *ComparePasswordCommand) Execute(ctx context.Context, username, password string) error {
	user, err := c.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return ErrPasswordMismatch
	}
	return nil
}
