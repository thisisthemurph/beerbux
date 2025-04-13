package builder

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"testing"
	"time"
)

type SessionBuilder struct {
	t                   *testing.T
	model               session.Session
	isActiveSetManually bool
	members             []SessionMemberParams
	existingMembers     []session.Member
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

func (b *SessionBuilder) WithUpdatedAt(t time.Time) *SessionBuilder {
	b.model.UpdatedAt = t
	return b
}

type SessionMemberParams struct {
	ID       string
	Name     string
	Username string
	IsOwner  bool
	IsAdmin  bool
}

func (b *SessionBuilder) WithMember(m SessionMemberParams) *SessionBuilder {
	b.members = append(b.members, m)
	return b
}

func (b *SessionBuilder) WithExistingMember(m session.Member) *SessionBuilder {
	b.existingMembers = append(b.existingMembers, m)
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
	insertSession := "insert into sessions (id, name, is_active, updated_at) values (?, ?, ?, ?);"
	insertMember := "insert into members (id, name, username) values (?, ?, ?);"
	insertSessionMember := "insert into session_members (session_id, member_id, is_owner, is_admin) values (?, ?, ?, ?);"
	insertTrans := "insert into transactions (id, session_id, member_id) values (?, ?, ?);"
	insertTransLine := "insert into transaction_lines (transaction_id, member_id, amount) values (?, ?, ?);"

	if !b.isActiveSetManually {
		b.model.IsActive = true
	}

	if b.model.ID == "" {
		b.model.ID = uuid.New().String()
	}

	_, err := db.Exec(insertSession, b.model.ID, b.model.Name, b.model.IsActive, b.model.UpdatedAt)
	require.NoError(b.t, err)

	for _, m := range b.members {
		_, err := db.Exec(insertMember, m.ID, m.Name, m.Username)
		require.NoError(b.t, err)

		_, err = db.Exec(insertSessionMember, b.model.ID, m.ID, m.IsOwner, m.IsAdmin)
		require.NoError(b.t, err)
	}

	for _, m := range b.existingMembers {
		_, err = db.Exec(insertSessionMember, b.model.ID, m.ID, false, false)
		require.NoError(b.t, err)
	}

	for _, t := range b.transactions {
		_, err := db.Exec(insertTrans, t.ID, t.SessionID, t.CreatorID)
		require.NoError(b.t, err)

		for _, line := range t.Lines {
			_, err := db.Exec(insertTransLine, t.ID, line.MemberID, line.Amount)
			require.NoError(b.t, err)
		}
	}

	return b.model
}
