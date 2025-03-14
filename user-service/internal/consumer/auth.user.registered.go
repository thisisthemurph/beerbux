package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
)

type UserRegisteredEvent struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserRegisteredKafkaConsumer struct {
	Reader         KafkaReader
	Logger         *slog.Logger
	UserRepository *user.Queries
}

func NewUserRegisteredKafkaConsumer(logger *slog.Logger, brokers []string, repo *user.Queries) *UserRegisteredKafkaConsumer {
	return &UserRegisteredKafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   "auth.user.registered",
			GroupID: "user-service",
		}),
		Logger:         logger,
		UserRepository: repo,
	}
}

func (c *UserRegisteredKafkaConsumer) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.Logger.Debug("Kafka consumer shutting down")
			return
		default:
			msg, err := c.Reader.ReadMessage(ctx)
			if ctx.Err() != nil || errors.Is(err, context.Canceled) {
				return
			}
			if err != nil {
				c.Logger.Error("Failed to read message", "error", err)
				continue
			}

			var ev UserRegisteredEvent
			if err := json.Unmarshal(msg.Value, &ev); err != nil {
				c.Logger.Error("Failed to unmarshal message", "error", err, "offset", msg.Offset)
				continue
			}

			_, err = c.UserRepository.CreateUser(ctx, user.CreateUserParams{
				ID:       ev.UserID,
				Name:     ev.Name,
				Username: ev.Username,
				Bio:      sql.NullString{},
			})
			if err != nil {
				c.Logger.Error("Failed to create user", "error", err)
				continue
			}
		}
	}
}

func (c *UserRegisteredKafkaConsumer) Close() error {
	return c.Reader.Close()
}
