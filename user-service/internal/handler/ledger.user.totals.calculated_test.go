package handler_test

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/user-service/internal/handler"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/tests/builder"
	"github.com/thisisthemurph/beerbux/user-service/tests/testinfra"
	"testing"
)

func TestLedgerUserTotalsCalculatedEventHandler_Handle(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	repo := user.New(db)
	h := handler.NewLedgerUserTotalsCalculatedEventHandler(repo)

	existingUser := builder.NewUserBuilder(t).
		WithName("Bob").
		WithUsername("bob").
		Build(db)

	event := handler.UserTotalsEvent{
		UserID: existingUser.ID,
		Credit: 100.0,
		Debit:  50.0,
		Net:    50.0,
	}

	eventBytes, err := json.Marshal(event)
	msg := kafka.Message{
		Key:   []byte(existingUser.ID),
		Value: eventBytes,
	}

	err = h.Handle(context.Background(), msg)
	require.NoError(t, err)

	q := "select credit, debit, net from users where id = ?;"
	var u user.User
	err = db.QueryRow(q, existingUser.ID).Scan(&u.Credit, &u.Debit, &u.Net)
	require.NoError(t, err)
	assert.Equal(t, event.Credit, u.Credit)
	assert.Equal(t, event.Debit, u.Debit)
	assert.Equal(t, event.Net, u.Net)
}
