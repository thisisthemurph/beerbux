package builder

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/pkg/nullish"
	"testing"
)

type UserBuilder struct {
	t     *testing.T
	model user.User
}

func NewUserBuilder(t *testing.T) *UserBuilder {
	return &UserBuilder{
		t: t,
	}
}

func (b *UserBuilder) WithID(id string) *UserBuilder {
	b.model.ID = id
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.model.Name = name
	return b
}

func (b *UserBuilder) WithUsername(username string) *UserBuilder {
	b.model.Username = username
	return b
}

func (b *UserBuilder) WithBio(bio string) *UserBuilder {
	b.model.Bio = nullish.CreateNullString(&bio)
	return b
}

func (b *UserBuilder) Build(db *sql.DB) user.User {
	if b.model.ID == "" {
		b.model.ID = uuid.NewString()
	}
	if b.model.Name == "" {
		b.t.Fatal("name is required")
	}
	if b.model.Username == "" {
		b.t.Fatal("username is required")
	}

	_, err := db.Exec(`
		insert into users (id, name, username, bio)
		values (?, ?, ?, ?);`,
		b.model.ID, b.model.Name, b.model.Username, b.model.Bio)

	if err != nil {
		b.t.Fatalf("failed to insert user: %v", err)
	}

	return b.model
}
