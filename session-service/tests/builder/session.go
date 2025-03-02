package builder

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
)

type SessionBuilder struct {
	t                   *testing.T
	model               session.Session
	isActiveSetManually bool
}

func NewSessionBuilder(t *testing.T) *SessionBuilder {
	return &SessionBuilder{
		t:     t,
		model: session.Session{},
	}
}

func (b *SessionBuilder) WithID(id uuid.UUID) *SessionBuilder {
	b.model.ID = id.String()
	return b
}

func (b *SessionBuilder) WithName(name string) *SessionBuilder {
	b.model.Name = name
	return b
}

func (b *SessionBuilder) WithIsActive(isActive bool) *SessionBuilder {
	b.isActiveSetManually = true
	b.model.IsActive = isActive
	return b
}

func (b *SessionBuilder) Build(db *sql.DB) session.Session {
	if !b.isActiveSetManually {
		b.model.IsActive = true
	}

	if b.model.ID == "" {
		b.model.ID = uuid.New().String()
	}

	_, err := db.Exec(`
		insert into sessions (id, name, is_active) 
		values (?, ?, ?);`,
		b.model.ID, b.model.Name, b.model.IsActive)

	if err != nil {
		b.t.Fatalf("failed to insert session: %v", err)
	}

	return b.model
}
