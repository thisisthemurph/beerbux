package consumer_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/user-service/internal/consumer"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/ledger"
	"github.com/thisisthemurph/beerbux/user-service/tests/testinfra"
)

type MockKafkaReader struct {
	mock.Mock
}

func (m *MockKafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (m *MockKafkaReader) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestLedgerUpdatedKafkaConsumer(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	repo := ledger.New(db)
	mockKafkaReader := new(MockKafkaReader)

	c := &consumer.LedgerUpdatedKafkaConsumer{
		Reader:               mockKafkaReader,
		Logger:               slog.Default(),
		UserLedgerRepository: repo,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		sessionID     = uuid.NewString()
		transactionID = uuid.NewString()
		creatorID     = uuid.NewString()
		participantID = uuid.NewString()
	)

	event1 := consumer.LedgerUpdatedEvent{
		ID:            uuid.NewString(),
		TransactionID: transactionID,
		SessionID:     sessionID,
		UserID:        creatorID,
		ParticipantID: participantID,
		Amount:        1,
	}

	event2 := consumer.LedgerUpdatedEvent{
		ID:            uuid.NewString(),
		TransactionID: transactionID,
		SessionID:     sessionID,
		UserID:        participantID,
		ParticipantID: creatorID,
		Amount:        -1,
	}

	allEvents := []consumer.LedgerUpdatedEvent{event1, event2}

	event1Bytes, _ := json.Marshal(event1)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{Value: event1Bytes}, nil).Once()
	event2Bytes, _ := json.Marshal(event2)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{Value: event2Bytes}, nil).Once()
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{}, context.Canceled)

	go c.Listen(ctx)
	time.Sleep(1 * time.Millisecond)

	mockKafkaReader.AssertExpectations(t)
	cancel()

	var res []ledger.UserLedger
	row, _ := db.Query("select id, user_id, participant_id, amount, type from user_ledger;")
	defer row.Close()
	for row.Next() {
		var r ledger.UserLedger
		err := row.Scan(&r.ID, &r.UserID, &r.ParticipantID, &r.Amount, &r.Type)
		require.NoError(t, err)
		res = append(res, r)
	}

	require.Len(t, res, len(allEvents))
	for i, r := range res {
		ev := allEvents[i]
		assert.Equal(t, ev.UserID, r.UserID)
		assert.Equal(t, ev.ParticipantID, r.ParticipantID)
		assert.Equal(t, ev.Amount, r.Amount)

		expectedType := "credit"
		if ev.Amount < 0 {
			expectedType = "debit"
		}
		assert.Equal(t, expectedType, r.Type)
	}
}
