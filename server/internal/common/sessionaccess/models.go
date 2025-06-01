package sessionaccess

import (
	"github.com/google/uuid"
	"time"
)

type SessionMember struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"isAdmin"`
	IsDeleted bool      `json:"isDeleted"`
}

type SessionTransactionLine struct {
	UserID uuid.UUID `json:"userId"`
	Amount float64   `json:"amount"`
}

type SessionTransaction struct {
	ID        uuid.UUID                `json:"id"`
	UserID    uuid.UUID                `json:"userId"`
	Total     float64                  `json:"total"`
	Lines     []SessionTransactionLine `json:"lines"`
	CreatedAt time.Time                `json:"createdAt"`
}

type SessionResponse struct {
	ID           uuid.UUID            `json:"id"`
	Name         string               `json:"name"`
	IsActive     bool                 `json:"isActive"`
	Members      []SessionMember      `json:"members"`
	Transactions []SessionTransaction `json:"transactions"`
	Total        float64              `json:"total"`
}
