package query

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

func (sr SessionResponse) GetMemberByID(userID uuid.UUID) (SessionMember, bool) {
	for _, m := range sr.Members {
		if m.ID == userID {
			return m, true
		}
	}
	return SessionMember{}, false
}

func (sr SessionResponse) IsMember(userID uuid.UUID) bool {
	m, ok := sr.GetMemberByID(userID)
	if ok {
		return !m.IsDeleted
	}
	return false
}

func (sr SessionResponse) IsAdminMember(userID uuid.UUID) bool {
	m, ok := sr.GetMemberByID(userID)
	if ok {
		return m.IsAdmin
	}
	return false
}
