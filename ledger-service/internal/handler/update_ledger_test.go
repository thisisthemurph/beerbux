package handler_test

import (
	"context"
	"database/sql"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/handler"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"github.com/thisisthemurph/beerbux/ledger-service/tests/testinfra"
)

// SetupTestHandler initializes the test database and handler
func SetupTestHandler(db *sql.DB) *handler.UpdateLedgerHandler {
	repo := repository.NewLedgerQueries(db)
	return handler.NewUpdateLedgerHandler(repo, slog.Default())
}

// CreateTestEvent generates a test event.TransactionCreatedEvent with given member amounts.
func CreateTestEvent(creatorID string, memberAmounts []event.TransactionCreatedMemberAmount) event.TransactionCreatedEvent {
	return event.TransactionCreatedEvent{
		TransactionID: uuid.NewString(),
		CreatorID:     creatorID,
		SessionID:     uuid.NewString(),
		MemberAmounts: memberAmounts,
	}
}

func TestHandle(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	h := SetupTestHandler(db)

	testCases := []struct {
		name          string
		memberAmounts []event.TransactionCreatedMemberAmount
		expectError   bool
	}{
		{
			name: "Valid transaction with two members",
			memberAmounts: []event.TransactionCreatedMemberAmount{
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 0.5},
			},
		},
		{
			name:          "No members (invalid case)",
			memberAmounts: []event.TransactionCreatedMemberAmount{},
			expectError:   true,
		},
		{
			name: "Single member (valid case)",
			memberAmounts: []event.TransactionCreatedMemberAmount{
				{UserID: uuid.NewString(), Amount: 1.5},
			},
		},
		{
			name: "Many members (valid case)",
			memberAmounts: []event.TransactionCreatedMemberAmount{
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
				{UserID: uuid.NewString(), Amount: 1.0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			creatorID := uuid.NewString()
			ev := CreateTestEvent(creatorID, tc.memberAmounts)

			res, err := h.Handle(context.Background(), ev)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
				return
			}

			assert.NoError(t, err, "Handler should not return an error")
			assert.Len(t, res, len(tc.memberAmounts)*2, "Should have double the entries (debits & credits)")

			// Validate debits and credits
			debits := make([]handler.MemberTransaction, 0)
			credits := make([]handler.MemberTransaction, 0)

			for _, r := range res {
				if r.UserID == creatorID {
					debits = append(debits, r)
				} else {
					credits = append(credits, r)
				}
			}

			assert.Len(t, debits, len(tc.memberAmounts), "Should have one debit per member")
			assert.Len(t, credits, len(tc.memberAmounts), "Should have one credit per member")

			for i, mt := range debits {
				assert.Equal(t, ev.TransactionID, mt.TransactionID, "Transaction ID mismatch")
				assert.Equal(t, ev.SessionID, mt.SessionID, "Session ID mismatch")
				assert.Equal(t, creatorID, mt.UserID, "UserID should match creator")
				assert.Equal(t, ev.MemberAmounts[i].UserID, mt.ParticipantID, "ParticipantID should match member")
				assert.Equal(t, -ev.MemberAmounts[i].Amount, mt.Amount, "Debit amount should be negative")
			}

			for i, mt := range credits {
				assert.Equal(t, ev.TransactionID, mt.TransactionID, "Transaction ID mismatch")
				assert.Equal(t, ev.SessionID, mt.SessionID, "Session ID mismatch")
				assert.Equal(t, creatorID, mt.ParticipantID, "ParticipantID should match creator")
				assert.Equal(t, ev.MemberAmounts[i].UserID, mt.UserID, "UserID should match member")
				assert.Equal(t, ev.MemberAmounts[i].Amount, mt.Amount, "Credit amount should match")
			}
		})
	}
}
