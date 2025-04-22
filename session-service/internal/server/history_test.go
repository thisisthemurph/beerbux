package server_test

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/history"
	"github.com/thisisthemurph/beerbux/session-service/internal/server"
	"github.com/thisisthemurph/beerbux/session-service/protos/historypb"
	"github.com/thisisthemurph/beerbux/session-service/tests/builder"
	"github.com/thisisthemurph/beerbux/session-service/tests/testinfra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"testing"
)

func TestGetBySessionID(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	historyRepo := history.NewHistoryRepository(db)
	historyServer := server.NewHistoryServer(historyRepo)

	sessionID := uuid.NewString()
	member1ID := uuid.NewString()
	member2ID := uuid.NewString()
	member3ID := uuid.NewString()

	eventData1 := &historypb.TransactionCreatedEventData{
		TransactionId: uuid.NewString(),
		Lines: []*historypb.TransactionLine{
			{
				MemberId: member2ID,
				Amount:   1,
			},
			{
				MemberId: member3ID,
				Amount:   1,
			},
		},
	}

	eventData2 := &historypb.TransactionCreatedEventData{
		TransactionId: uuid.NewString(),
		Lines: []*historypb.TransactionLine{
			{
				MemberId: member1ID,
				Amount:   1,
			},
			{
				MemberId: member3ID,
				Amount:   1,
			},
		},
	}

	data1, err := json.Marshal(eventData1)
	require.NoError(t, err)
	data2, err := json.Marshal(eventData2)
	require.NoError(t, err)

	ev1 := builder.NewSessionHistoryBuilder(t).
		WithID(1).
		WithSessionID(sessionID).
		WithMemberID(member1ID).
		WithEventType(history.EventTransactionCreated).
		WithEventData(data1).
		Build(db)
	ev2 := builder.NewSessionHistoryBuilder(t).
		WithID(2).
		WithSessionID(sessionID).
		WithMemberID(member2ID).
		WithEventType(history.EventTransactionCreated).
		WithEventData(data2).
		Build(db)
	// A session history item for a different session that should not be retrieved.
	_ = builder.NewSessionHistoryBuilder(t).
		WithID(3).
		WithSessionID(uuid.NewString()).
		WithMemberID(uuid.NewString()).
		WithEventType(history.EventTransactionCreated).
		WithEventData(data2).
		Build(db)

	expectedEvents := []history.SessionHistory{ev1, ev2}

	sessionHistory, err := historyServer.GetBySessionID(context.Background(), &historypb.GetBySessionIDRequest{
		SessionId: sessionID,
	})

	assert.NoError(t, err)
	assert.NotNil(t, sessionHistory)
	assert.Equal(t, sessionID, sessionHistory.SessionId)
	assert.Len(t, sessionHistory.Events, len(expectedEvents))

	for i, event := range sessionHistory.Events {
		expectedEvent := expectedEvents[i]
		assert.Equal(t, expectedEvent.ID, event.Id)
		assert.Equal(t, expectedEvent.MemberID, event.MemberId)
		assert.Equal(t, expectedEvent.EventType, event.EventType)
		assert.NotNil(t, event.EventData)
		assertEventData(t, event.EventType, expectedEvent.EventData, event.EventData)
	}
}

func assertEventData(t *testing.T, eventType string, expectedData []byte, data *anypb.Any) {
	et := history.NewEventType(eventType)
	switch et {
	case history.EventTransactionCreated:
		var expectedTransactionCreatedEventData *historypb.TransactionCreatedEventData
		err := json.Unmarshal(expectedData, &expectedTransactionCreatedEventData)
		require.NoError(t, err)
		assertTransactionCreatedEventData(t, expectedTransactionCreatedEventData, data)
	default:
		t.Fatalf("unknown event type %s", eventType)
	}
}

func assertTransactionCreatedEventData(t *testing.T, expected *historypb.TransactionCreatedEventData, data *anypb.Any) {
	var eventData historypb.TransactionCreatedEventData
	err := anypb.UnmarshalTo(data, &eventData, proto.UnmarshalOptions{})
	require.NoError(t, err)

	assert.Equal(t, expected.TransactionId, eventData.TransactionId)
	assert.Equal(t, len(expected.Lines), len(eventData.Lines))
	for i, line := range eventData.Lines {
		expectedLine := expected.Lines[i]
		assert.Equal(t, expectedLine.MemberId, line.MemberId)
		assert.Equal(t, expectedLine.Amount, line.Amount)
	}
}
