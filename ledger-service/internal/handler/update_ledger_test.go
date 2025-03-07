package handler_test

import (
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/handler"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
	"github.com/thisisthemurph/beerbux/ledger-service/tests/testinfra"
)

func TestHandle_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	repo := repository.NewLedgerQueries(db)
	h := handler.NewUpdateLedgerHandler(repo, slog.Default())

	ev := event.TransactionCreatedEvent{
		Metadata: event.Metadata{
			Version: "1.0.0",
		},
		Data: event.TransactionCreatedEventData{
			TransactionID: uuid.New(),
			CreatorID:     uuid.New(),
			SessionID:     uuid.New(),
			MemberAmounts: []event.TransactionCreatedMemberAmount{
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

	res, err := h.Handle(ev)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	for i, r := range res {
		assert.NotEmpty(t, r.ID)
		assert.Equal(t, ev.Data.TransactionID, r.TransactionID)
		assert.Equal(t, ev.Data.SessionID.String(), r.SessionID.String())
		assert.Equal(t, ev.Data.MemberAmounts[i].UserID.String(), r.UserID.String())
		assert.Equal(t, ev.Data.MemberAmounts[i].Amount, r.Amount)
	}

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
