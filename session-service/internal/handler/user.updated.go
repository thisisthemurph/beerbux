package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
)

type UserUpdatedEventHandler struct {
	sessionRepository *session.Queries
	logger            *slog.Logger
}

func NewUserUpdatedEventHandler(sessionRepository *session.Queries, logger *slog.Logger) *UserUpdatedEventHandler {
	return &UserUpdatedEventHandler{
		sessionRepository: sessionRepository,
		logger:            logger,
	}
}

type UserUpdatedEventData struct {
	UserID        string                 `json:"user_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"`
}

type UserUpdatedEvent struct {
	Data UserUpdatedEventData `json:"user"`
}

func (h *UserUpdatedEventHandler) Handle(msg *nats.Msg) {
	var event UserUpdatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		h.logger.Error("failed to unmarshal user.updated event", "error", err)
		return
	}

	existingMember, err := h.getExistingMember(event.Data.UserID)
	if err != nil {
		return
	}

	err = h.sessionRepository.UpdateMember(context.TODO(), session.UpdateMemberParams{
		ID:       event.Data.UserID,
		Name:     stringFieldOrDefault(event.Data.UpdatedFields, "name", existingMember.Name),
		Username: stringFieldOrDefault(event.Data.UpdatedFields, "username", existingMember.Username),
	})

	if err != nil {
		h.logger.Error("failed to update user", "error", err)
		return
	}
}

func (h *UserUpdatedEventHandler) getExistingMember(userID string) (session.Member, error) {
	existingMember, err := h.sessionRepository.GetMember(context.TODO(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Debug("user not found", "user_id", userID)
		} else {
			h.logger.Error("failed to get user", "error", err)
		}
		return session.Member{}, err
	}

	return existingMember, nil
}

// stringFieldOrDefault returns the value of the key in the data map if present and a string,
// otherwise it returns the default value.
func stringFieldOrDefault(data map[string]interface{}, key, defaultValue string) string {
	value, ok := data[key].(string)
	if !ok {
		return defaultValue
	}
	return value
}
