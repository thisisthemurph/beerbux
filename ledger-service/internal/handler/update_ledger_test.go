package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/handler"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"github.com/thisisthemurph/beerbux/ledger-service/tests/testinfra"
)

// NewTestUpdateLedgerHandler initializes the test database and handler.
func NewTestUpdateLedgerHandler(db *sql.DB) (*handler.UpdateLedgerHandler, chan []event.LedgerUpdateEvent) {
	repo := repository.NewLedgerQueries(db)
	c := make(chan []event.LedgerUpdateEvent, 10)
	return handler.NewUpdateLedgerHandler(repo, c, slog.Default()), c
}

// CreateKafkaMessageAndEvent generates a test event.TransactionCreatedEvent with given member amounts.
func CreateKafkaMessageAndEvent(
	t *testing.T,
	creatorID string,
	memberAmounts []event.TransactionCreatedMemberAmount,
) (kafka.Message, event.TransactionCreatedEvent) {
	baseEvent := event.TransactionCreatedEvent{
		TransactionID: uuid.NewString(),
		CreatorID:     creatorID,
		SessionID:     uuid.NewString(),
		MemberAmounts: memberAmounts,
	}

	data, err := json.Marshal(baseEvent)
	require.NoError(t, err)

	return kafka.Message{
		Key:   []byte(baseEvent.TransactionID),
		Value: data,
	}, baseEvent
}

func TestHandle(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	h, resultChan := NewTestUpdateLedgerHandler(db)

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
			msg, ev := CreateKafkaMessageAndEvent(t, creatorID, tc.memberAmounts)

			err := h.Handle(context.Background(), msg)

			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
				return
			} else {
				require.NoError(t, err, "Handler should not return an error")
				if err != nil {
					return
				}
			}

			res := <-resultChan
			assert.Len(t, res, len(tc.memberAmounts)*2, "Should have double the entries (debits & credits)")

			// Validate debits and credits
			debits := make([]event.LedgerUpdateEvent, 0)
			credits := make([]event.LedgerUpdateEvent, 0)

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
