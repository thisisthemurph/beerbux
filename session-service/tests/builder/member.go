package builder

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
)

type MemberBuilder struct {
	t     *testing.T
	model session.Member
}

func NewMemberBuilder(t *testing.T) *MemberBuilder {
	return &MemberBuilder{t: t}
}

func (b *MemberBuilder) WithID(id uuid.UUID) *MemberBuilder {
	b.model.ID = id.String()
	return b
}

func (b *MemberBuilder) WithName(name string) *MemberBuilder {
	b.model.Name = name
	return b
}

func (b *MemberBuilder) WithUsername(username string) *MemberBuilder {
	b.model.Username = username
	return b
}

func (b *MemberBuilder) Build(db *sql.DB) session.Member {
	if b.model.ID == "" {
		b.model.ID = uuid.NewString()
	}
	if b.model.Name == "" {
		b.t.Fatalf("name is required")
	}
	if b.model.Username == "" {
		b.t.Fatalf("username is required")
	}

	_, err := db.Exec(
		"insert into members (id, name, username) values (?, ?, ?)",
		b.model.ID, b.model.Name, b.model.Username)

	if err != nil {
		b.t.Fatalf("failed to insert member: %v", err)
	}

	return b.model
}
