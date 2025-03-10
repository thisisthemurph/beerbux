package handler_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/beerbux/session-service/internal/handler"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/tests/builder"
	"github.com/thisisthemurph/beerbux/session-service/tests/testinfra"
	_ "modernc.org/sqlite"
)

func TestUserUpdatedHandler_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	h := handler.NewUserUpdatedEventHandler(sessionRepo)

	member := builder.NewMemberBuilder(t).
		WithName("John Doe").
		WithUsername("johndoe").
		Build(db)

	data := map[string]interface{}{
		"user_id": member.ID,
		"updated_fields": map[string]interface{}{
			"name":       "Updated Name",
			"username":   "updated.name",
			"updated_at": "2021-09-01T00:00:00Z",
		},
	}

	msgData, err := json.Marshal(data)
	assert.NoError(t, err)

	err = h.Handle(kafka.Message{Value: msgData})
	assert.NoError(t, err)

	var name, username string
	var updatedAt time.Time
	err = db.QueryRow("select name, username, updated_at from members where id = ?", member.ID).Scan(&name, &username, &updatedAt)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", name)
	assert.Equal(t, "updated.name", username)
	assert.WithinDuration(t, time.Now(), updatedAt, 5*time.Second)
	assert.Greater(t, updatedAt.Unix(), member.UpdatedAt.Unix())
}
