package handler

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository"
	"github.com/thisisthemurph/beerbux/session-service/tests/testinfra"
)

func TestLedgerTransactionUpdatedMessageHandler_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := repository.NewSessionQueries(db)
	handler := NewTransactionCreatedMessageHandler(sessionRepo)

	event := TransactionCreatedEvent{
		TransactionID: uuid.NewString(),
		SessionID:     uuid.NewString(),
		CreatorID:     uuid.NewString(),
		Amounts: []MemberAmount{
			{
				MemberID: "a3d7ec98-034d-4c25-8b9e-023faa19fd37",
				Amount:   1,
			},
			{
				MemberID: "f3d7ec98-034d-4c25-8b9e-023faa19fd37",
				Amount:   1,
			},
		},
	}

	msgData, _ := json.Marshal(event)
	msg := kafka.Message{Value: msgData}

	err := handler.Handle(context.Background(), msg)
	assert.NoError(t, err)

	q := "select id, session_id, member_id from transactions where id = ?"
	var transactionID, sessionID, memberID string
	err = db.QueryRow(q, event.TransactionID).Scan(&transactionID, &sessionID, &memberID)
	assert.NoError(t, err)
	assert.Equal(t, event.TransactionID, transactionID)
	assert.Equal(t, event.SessionID, sessionID)
	assert.Equal(t, event.CreatorID, memberID)

	q = "select member_id, amount from transaction_lines where transaction_id = ? order by member_id"
	var amount float64
	rows, err := db.Query(q, event.TransactionID)
	assert.NoError(t, err)
	defer rows.Close()
	var index int
	for rows.Next() {
		err = rows.Scan(&memberID, &amount)
		assert.NoError(t, err)
		assert.Equal(t, 1.0, amount)
		assert.Equal(t, event.Amounts[index].MemberID, memberID)
		index++
	}
}
