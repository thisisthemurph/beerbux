package useraccess

import (
	"context"
	"github.com/google/uuid"
)

type UserReader interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
}
