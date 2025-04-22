package builder

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/history"
	"testing"
)

type SessionHistoryBuilder struct {
	t     *testing.T
	model history.SessionHistory
}

func NewSessionHistoryBuilder(t *testing.T) *SessionHistoryBuilder {
	return &SessionHistoryBuilder{
		t:     t,
		model: history.SessionHistory{},
	}
}

func (b *SessionHistoryBuilder) WithID(id int64) *SessionHistoryBuilder {
	b.model.ID = id
	return b
}

func (b *SessionHistoryBuilder) WithSessionID(id string) *SessionHistoryBuilder {
	b.model.SessionID = id
	return b
}

func (b *SessionHistoryBuilder) WithMemberID(id string) *SessionHistoryBuilder {
	b.model.MemberID = id
	return b
}

func (b *SessionHistoryBuilder) WithEventType(et history.EventType) *SessionHistoryBuilder {
	b.model.EventType = et.String()
	return b
}

func (b *SessionHistoryBuilder) WithEventData(data []byte) *SessionHistoryBuilder {
	b.model.EventData = data
	return b
}

func (b *SessionHistoryBuilder) Build(db *sql.DB) history.SessionHistory {
	q := "insert into session_history (id, session_id, member_id, event_type, event_data) values (?, ?, ?, ?, ?);"
	_, err := db.Exec(q, b.model.ID, b.model.SessionID, b.model.MemberID, b.model.EventType, b.model.EventData)
	require.NoError(b.t, err)
	return b.model
}
