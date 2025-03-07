package handler_test

import (
	"encoding/json"
	"log/slog"
	"testing"
	
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/handler"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"github.com/thisisthemurph/beerbux/ledger-service/tests/testinfra"
)

func TestHandle_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	repo := repository.NewLedgerQueries(db)
	h := handler.NewTransactionCreatedMsgHandler(repo, slog.Default())

	ev := handler.TransactionCreatedEvent{
		EventMetadata: handler.EventMetadata{
			Version: "1.0.0",
		},
		Data: handler.TransactionCreatedEventData{
			TransactionID: uuid.New(),
			CreatorID:     uuid.New(),
			SessionID:     uuid.New(),
			MemberAmounts: []handler.TransactionCreatedMemberAmount{
				{
					UserID: uuid.New(),
					Amount: 1,
				},
				{
					UserID: uuid.New(),
					Amount: 0.5,
				},
			},
		},
	}

	eventData, err := json.Marshal(ev)
	assert.NoError(t, err)
	msg := nats.Msg{Data: eventData}

	h.Handle(&msg)

	query := `select session_id, user_id, amount from ledger where transaction_id = ?;`
	rows, err := db.Query(query, ev.Data.TransactionID)
	assert.NoError(t, err)
	defer rows.Close()

	var sessionID, userID string
	var amount float64

	i := 0
	for rows.Next() {
		err := rows.Scan(&sessionID, &userID, &amount)
		assert.NoError(t, err)
		assert.Equal(t, ev.Data.SessionID.String(), sessionID)
		assert.Equal(t, ev.Data.MemberAmounts[i].UserID.String(), userID)
		assert.Equal(t, ev.Data.MemberAmounts[i].Amount, amount)
		i++
	}

	assert.Equal(t, len(ev.Data.MemberAmounts), i)
}
