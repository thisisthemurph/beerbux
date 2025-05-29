package useraccess

import (
	"context"
	"github.com/google/uuid"
)

type UserReader interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	GetUserByUsername(ctx context.Context, username string) (*UserResponse, error)
	GetUserByEmail(ctx context.Context, username string) (*UserResponse, error)
	UserWithUsernameExists(ctx context.Context, username string) (bool, error)
	UserWithEmailExists(ctx context.Context, username string) (bool, error)
}
