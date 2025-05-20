package useraccess

import (
	"context"
	"github.com/google/uuid"
)

type UserReader interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	GetUserByUsername(ctx context.Context, username string) (*UserResponse, error)
}
