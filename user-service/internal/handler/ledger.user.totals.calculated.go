package handler

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
)

type LedgerUserTotalsCalculatedEventHandler struct {
	userRepository *user.Queries
}

func NewLedgerUserTotalsCalculatedEventHandler(userRepository *user.Queries) KafkaMessageHandler {
	return &LedgerUserTotalsCalculatedEventHandler{
		userRepository: userRepository,
	}
}

type UserTotalsEvent struct {
	UserID string  `json:"user_id"`
	Credit float64 `json:"credit"`
	Debit  float64 `json:"debit"`
	Net    float64 `json:"net"`
}

func (h *LedgerUserTotalsCalculatedEventHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var userTotalsEvent UserTotalsEvent
	if err := json.Unmarshal(msg.Value, &userTotalsEvent); err != nil {
		return err
	}

	err := h.userRepository.UpdateUserTotals(ctx, user.UpdateUserTotalsParams{
		ID:     userTotalsEvent.UserID,
		Credit: userTotalsEvent.Credit,
		Debit:  userTotalsEvent.Debit,
		Net:    userTotalsEvent.Net,
	})

	return err
}
