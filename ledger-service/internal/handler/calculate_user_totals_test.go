package handler_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/handler"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"github.com/thisisthemurph/beerbux/ledger-service/tests/testinfra"
)

func NewTestCalculateUserTotalsHandler(db *sql.DB) *handler.CalculateUserTotalsHandler {
	repo := repository.NewLedgerQueries(db)
	return handler.NewCalculateUserTotalsHandler(repo)
}

func TestCalculateUserTotals(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	h := NewTestCalculateUserTotalsHandler(db)

	q := `	
		insert into ledger (id, transaction_id, session_id, user_id, amount)
		values (?, ?, ?, ?, ?);`

	user1ID := uuid.NewString()
	user2ID := uuid.NewString()
	user3ID := uuid.NewString()
	nonExistentUserID := uuid.NewString()

	// Ledger items for user1
	_, _ = db.Exec(q, uuid.NewString(), uuid.NewString(), uuid.NewString(), user1ID, 1)
	_, _ = db.Exec(q, uuid.NewString(), uuid.NewString(), uuid.NewString(), user1ID, 1)
	_, _ = db.Exec(q, uuid.NewString(), uuid.NewString(), uuid.NewString(), user1ID, 1)
	_, _ = db.Exec(q, uuid.NewString(), uuid.NewString(), uuid.NewString(), user1ID, -1)
	_, _ = db.Exec(q, uuid.NewString(), uuid.NewString(), uuid.NewString(), user1ID, -1)

	// Ledger items for user2 (only debits)
	_, _ = db.Exec(q, uuid.NewString(), uuid.NewString(), uuid.NewString(), user2ID, -5)
	_, _ = db.Exec(q, uuid.NewString(), uuid.NewString(), uuid.NewString(), user2ID, -2)

	// Ledger items for user3 (only credits)
	_, _ = db.Exec(q, uuid.NewString(), uuid.NewString(), uuid.NewString(), user3ID, 3)
	_, _ = db.Exec(q, uuid.NewString(), uuid.NewString(), uuid.NewString(), user3ID, 4)

	// Test user1
	res, err := h.Handle(context.Background(), user1ID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user1ID, res.UserID)
	assert.Equal(t, 2.0, res.Credit)
	assert.Equal(t, 3.0, res.Debit)
	assert.Equal(t, 1.0, res.Net)

	// Test user2 (only debits)
	res, err = h.Handle(context.Background(), user2ID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user2ID, res.UserID)
	assert.Equal(t, 7.0, res.Credit)
	assert.Equal(t, 0.0, res.Debit)
	assert.Equal(t, -7.0, res.Net)

	// Test user3 (only credits)
	res, err = h.Handle(context.Background(), user3ID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user3ID, res.UserID)
	assert.Equal(t, 0.0, res.Credit)
	assert.Equal(t, 7.0, res.Debit)
	assert.Equal(t, 7.0, res.Net)

	// Test non-existent user
	res, err = h.Handle(context.Background(), nonExistentUserID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, nonExistentUserID, res.UserID)
	assert.Equal(t, 0.0, res.Credit)
	assert.Equal(t, 0.0, res.Debit)
	assert.Equal(t, 0.0, res.Net)
}
