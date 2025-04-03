package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
)

type AuthUserRegisteredHandler struct {
	userRepository *user.Queries
}

func NewAuthUserRegisteredHandler(userRepository *user.Queries) KafkaMessageHandler {
	return &AuthUserRegisteredHandler{
		userRepository: userRepository,
	}
}

type UserRegisteredEvent struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (c *AuthUserRegisteredHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var ev UserRegisteredEvent
	if err := json.Unmarshal(msg.Value, &ev); err != nil {
		return err
	}

	_, err := c.userRepository.CreateUser(ctx, user.CreateUserParams{
		ID:       ev.UserID,
		Name:     ev.Name,
		Username: ev.Username,
		Bio:      sql.NullString{},
	})

	if err != nil {
		return fmt.Errorf("failed to create user with id %s: %w", ev.UserID, err)
	}

	return nil
}
