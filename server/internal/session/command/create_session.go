package command

import (
	"beerbux/internal/common/useraccess"
	"beerbux/internal/session/db"
	"beerbux/pkg/dbtx"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type CreateSessionCommand struct {
	dbtx.TX
	Queries    *db.Queries
	UserReader useraccess.UserReader
}

func NewCreateSessionCommand(queries *db.Queries, userReader useraccess.UserReader) *CreateSessionCommand {
	return &CreateSessionCommand{
		Queries:    queries,
		UserReader: userReader,
	}
}

type SessionMember struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	IsCreator bool      `json:"isCreator"`
	IsAdmin   bool      `json:"isAdmin"`
	IsDeleted bool      `json:"isDeleted"`
}

type CreateSessionResponse struct {
	ID       uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	IsActive bool            `json:"isActive"`
	Members  []SessionMember `json:"members"`
}

func (c *CreateSessionCommand) Execute(ctx context.Context, userID uuid.UUID, sessionName string) (*CreateSessionResponse, error) {
	if err := validateSessionName(sessionName); err != nil {
		return nil, fmt.Errorf("invalid sessin name: %w", err)
	}

	user, err := c.UserReader.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, useraccess.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user with id %s: %w", userID, err)
	}

	tx, err := c.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	qtx := c.Queries.WithTx(tx)

	session, err := qtx.CreateSession(ctx, db.CreateSessionParams{
		Name:      sessionName,
		CreatorID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	err = qtx.AddMemberToSession(ctx, db.AddMemberToSessionParams{
		SessionID: session.ID,
		MemberID:  user.ID,
		IsAdmin:   true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add member to session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	return &CreateSessionResponse{
		ID:       session.ID,
		Name:     session.Name,
		IsActive: true,
		Members: []SessionMember{
			{
				ID:        user.ID,
				Name:      user.Name,
				Username:  user.Username,
				IsCreator: true,
				IsAdmin:   true,
			},
		},
	}, nil
}

func validateSessionName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) < 3 {
		return fmt.Errorf("name must be at least 3 characters")
	}
	if len(name) > 20 {
		return fmt.Errorf("name must be at most 20 characters")
	}
	return nil
}
