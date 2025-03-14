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
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/tests/fake"
	"github.com/thisisthemurph/beerbux/user-service/tests/testinfra"
)

func TestUserRegisteredKafkaConsumer(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	repo := user.New(db)
	mockKafkaReader := new(fake.MockKafkaReader)

	c := &consumer.UserRegisteredKafkaConsumer{
		Reader:         mockKafkaReader,
		Logger:         slog.Default(),
		UserRepository: repo,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	event1 := consumer.UserRegisteredEvent{
		UserID:   uuid.NewString(),
		Name:     "Alice",
		Username: "wonderful.alice",
	}

	event2 := consumer.UserRegisteredEvent{
		UserID:   uuid.NewString(),
		Name:     "Bob",
		Username: "bobs.coffee",
	}

	allEvents := []consumer.UserRegisteredEvent{event1, event2}

	eventBytes, _ := json.Marshal(event1)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{Value: eventBytes}, nil).Once()
	eventBytes2, _ := json.Marshal(event2)
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{Value: eventBytes2}, nil).Once()
	mockKafkaReader.On("ReadMessage", mock.Anything).Return(kafka.Message{}, context.Canceled)

	go c.Listen(ctx)
	time.Sleep(2 * time.Millisecond)

	mockKafkaReader.AssertExpectations(t)
	cancel()

	var userRecord user.User
	row, err := db.Query("SELECT id, name, username FROM users order by name;", event1.UserID)
	require.NoError(t, err)
	defer row.Close()

	var i int
	for row.Next() {
		expectedUser := allEvents[i]
		err = row.Scan(&userRecord.ID, &userRecord.Name, &userRecord.Username)
		require.NoError(t, err)

		assert.Equal(t, expectedUser.UserID, userRecord.ID)
		assert.Equal(t, expectedUser.Name, userRecord.Name)
		assert.Equal(t, expectedUser.Username, userRecord.Username)
		i++
	}
}
