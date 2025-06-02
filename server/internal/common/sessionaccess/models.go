package sessionaccess

import (
	"github.com/google/uuid"
	"time"
)

// Session contains the basic details of a session, including a slice of members.
type Session struct {
	ID       uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	IsActive bool            `json:"isActive"`
	Members  []SessionMember `json:"members"`
}

// SessionWithTransactions contains all Session data, plus the SessionTransactions slice and a Total indicating the total of all transactions.
type SessionWithTransactions struct {
	Session
	Transactions []SessionTransaction `json:"transactions"`
	Total        float64              `json:"total"`
}

// HasMember returns a bool indicating if the session contains a member of the given ID.
func (s *Session) HasMember(memberID uuid.UUID) bool {
	for _, m := range s.Members {
		if m.ID == memberID {
			return true
		}
	}
	return false
}

// HasAdminMember returns a bool indicating if the session has an admin member of the given ID.
// Returns false if there is no such member, or if there is a member and the member is not an admin.
func (s *Session) HasAdminMember(memberID uuid.UUID) bool {
	for _, m := range s.Members {
		if m.ID == memberID {
			return m.IsAdmin
		}
	}
	return false
}

type SessionMember struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"isAdmin"`
	IsDeleted bool      `json:"isDeleted"`
}

// SessionTransaction is a transaction initiated by one member to one or more other members within the session.
type SessionTransaction struct {
	ID        uuid.UUID                `json:"id"`
	UserID    uuid.UUID                `json:"userId"`
	Total     float64                  `json:"total"`
	Lines     []SessionTransactionLine `json:"lines"`
	CreatedAt time.Time                `json:"createdAt"`
}

// SessionTransactionLine is a component of the SessionTransaction; a smaller part of a transaction to the UserID.
type SessionTransactionLine struct {
	UserID uuid.UUID `json:"userId"`
	Amount float64   `json:"amount"`
}
