package builder

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"testing"
)

type SessionBuilder struct {
	t                   *testing.T
	model               session.Session
	isActiveSetManually bool
	members             []SessionMemberParams
	transactions        []SessionTransactionParams
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

type SessionMemberParams struct {
	ID       string
	Name     string
	Username string
	IsOwner  bool
}

func (b *SessionBuilder) WithMember(m SessionMemberParams) *SessionBuilder {
	b.members = append(b.members, m)
	return b
}

type SessionTransactionLine struct {
	MemberID string
	Amount   float64
}

type SessionTransactionParams struct {
	ID        string
	SessionID string
	CreatorID string
	Lines     []SessionTransactionLine
}

func (b *SessionBuilder) WithTransaction(t SessionTransactionParams) *SessionBuilder {
	b.transactions = append(b.transactions, t)
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

	for _, m := range b.members {
		_, err := db.Exec(`
			insert into members (id, name, username) 
			values (?, ?, ?);`,
			m.ID, m.Name, m.Username)

		if err != nil {
			b.t.Fatalf("failed to insert member: %v", err)
		}

		_, err = db.Exec(`
			insert into session_members (session_id, member_id, is_owner) 
			values (?, ?, ?);`,
			b.model.ID, m.ID, m.IsOwner)

		if err != nil {
			b.t.Fatalf("failed to insert session member: %v", err)
		}
	}

	insertTransaction := "insert into transactions (id, session_id, member_id) values (?, ?, ?);"
	insertTransactionLine := "insert into transaction_lines (transaction_id, member_id, amount) values (?, ?, ?);"
	for _, t := range b.transactions {
		_, err := db.Exec(insertTransaction, t.ID, t.SessionID, t.CreatorID)
		if err != nil {
			b.t.Fatalf("failed to insert transaction: %v", err)
		}
		for _, line := range t.Lines {
			_, err := db.Exec(insertTransactionLine, t.ID, line.MemberID, line.Amount)
			if err != nil {
				b.t.Fatalf("failed to insert transaction line: %v", err)
			}
		}
	}

	return b.model
}
