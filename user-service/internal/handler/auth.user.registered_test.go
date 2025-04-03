package handler_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/user-service/internal/handler"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/tests/testinfra"
)

func TestUserRegisteredHandler(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	repo := user.New(db)
	h := handler.NewAuthUserRegisteredHandler(repo)

	event := handler.UserRegisteredEvent{
		UserID:   "12345",
		Name:     "Alice",
		Username: "alice.in.wonderland",
	}

	eventBytes, err := json.Marshal(event)
	require.NoError(t, err)

	msg := kafka.Message{
		Key:   []byte(event.UserID),
		Value: eventBytes,
	}

	err = h.Handle(context.Background(), msg)
	require.NoError(t, err)

	q := "select name, username, balance from users where id = ?;"
	var u user.User
	err = db.QueryRow(q, event.UserID).Scan(&u.Name, &u.Username, &u.Balance)
	require.NoError(t, err)
	assert.Equal(t, event.Name, u.Name)
	assert.Equal(t, event.Username, u.Username)
	assert.Equal(t, 0.0, u.Balance)
}
