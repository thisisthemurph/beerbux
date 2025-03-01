package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/session-service/internal/events"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/user"
)

type UserCreatedMsgHandler struct {
	userRepository *user.Queries
	logger         *slog.Logger
}

func NewUserCreatedMsgHandler(userRepository *user.Queries) *UserCreatedMsgHandler {
	return &UserCreatedMsgHandler{
		userRepository: userRepository,
	}
}

func (h *UserCreatedMsgHandler) Handle(m *nats.Msg) {
	var ev events.UserCreatedEvent
	if err := json.Unmarshal(m.Data, &ev); err != nil {
		h.logger.Error("failed to unmarshal user.created message", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.userRepository.CreateUser(ctx, user.CreateUserParams{
		ID:       ev.User.ID,
		Username: ev.User.Username,
	}); err != nil {
		h.logger.Error("failed to create user", "error", err)
	}
}
