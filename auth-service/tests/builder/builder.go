package builder

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/auth"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type UserBuilder struct {
	t     *testing.T
	model auth.User
}

func NewUserBuilder(t *testing.T) *UserBuilder {
	return &UserBuilder{
		t: t,
	}
}

func (b *UserBuilder) WithID(id uuid.UUID) *UserBuilder {
	b.model.ID = id.String()
	return b
}

func (b *UserBuilder) WithUsername(username string) *UserBuilder {
	b.model.Username = username
	return b
}

// WithPassword takes the raw string password and hashes it for the user.hashed_password column.
func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		b.t.Fatalf("failed to hash password: %v", err)
	}
	b.model.HashedPassword = string(hashedBytes)
	return b
}

func (b *UserBuilder) Build(db *sql.DB) auth.User {
	if b.model.ID == "" {
		b.model.ID = uuid.NewString()
	}

	stmt := "insert into users (id, username, hashed_password) values (?, ?, ?);"
	_, err := db.Exec(stmt, b.model.ID, b.model.Username, b.model.HashedPassword)
	if err != nil {
		b.t.Fatalf("failed to insert user: %v", err)
	}

	b.t.Log(fmt.Sprintf("builder: user created: %v (%v)", b.model.Username, b.model.ID))
	return b.model
}
