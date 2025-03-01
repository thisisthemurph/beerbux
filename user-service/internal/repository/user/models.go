// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string
	Name      string
	Username  string
	Bio       sql.NullString
	Balance   float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
