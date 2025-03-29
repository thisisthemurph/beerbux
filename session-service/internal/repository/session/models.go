// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package session

import (
	"time"
)

type Member struct {
	ID        string
	Name      string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID        string
	Name      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SessionMember struct {
	SessionID string
	MemberID  string
	IsOwner   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID        string
	SessionID string
	MemberID  string
}

type TransactionLine struct {
	TransactionID string
	MemberID      string
	Amount        float64
}
