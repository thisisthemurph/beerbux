package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
)

type UserUpdatedEventHandler struct {
	sessionRepository *session.Queries
}

func NewUserUpdatedEventHandler(sessionRepository *session.Queries) *UserUpdatedEventHandler {
	return &UserUpdatedEventHandler{
		sessionRepository: sessionRepository,
	}
}

type UserUpdatedEvent struct {
	UserID        string                 `json:"user_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"`
}

func (h *UserUpdatedEventHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var updatedUser UserUpdatedEvent
	if err := json.Unmarshal(msg.Value, &updatedUser); err != nil {
		return err
	}

	existingMember, err := h.getExistingMember(ctx, updatedUser.UserID)
	if err != nil {
		return err
	}

	err = h.sessionRepository.UpdateMember(context.TODO(), session.UpdateMemberParams{
		ID:       updatedUser.UserID,
		Name:     stringFieldOrDefault(updatedUser.UpdatedFields, "name", existingMember.Name),
		Username: stringFieldOrDefault(updatedUser.UpdatedFields, "username", existingMember.Username),
	})

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// getExistingMember retrieves the existing member from the database.
func (h *UserUpdatedEventHandler) getExistingMember(ctx context.Context, userID string) (session.Member, error) {
	existingMember, err := h.sessionRepository.GetMember(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return session.Member{}, fmt.Errorf("user not found: %w", err)
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
