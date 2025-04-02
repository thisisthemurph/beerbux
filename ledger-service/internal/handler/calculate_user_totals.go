package handler

import (
	"context"
	"database/sql"
	"errors"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/event"
	"github.com/thisisthemurph/beerbux/ledger-service/internal/repository"
)

type CalculateUserTotalsHandler struct {
	ledgerRepository *repository.LedgerQueriesWrapper
}

func NewCalculateUserTotalsHandler(ledgerRepository *repository.LedgerQueriesWrapper) *CalculateUserTotalsHandler {
	return &CalculateUserTotalsHandler{
		ledgerRepository: ledgerRepository,
	}
}

func (h *CalculateUserTotalsHandler) Handle(ctx context.Context, userID string) (event.UserTotalsEvent, error) {
	res, err := h.ledgerRepository.CalculateUserTotal(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If the user has no associated ledger items, we will return zero balance.
			return event.UserTotalsEvent{
				UserID: userID,
			}, nil
		}
		return event.UserTotalsEvent{}, err
	}

	return event.UserTotalsEvent{
		UserID: res.UserID,
		Credit: res.Credit,
		Debit:  res.Debit,
		Net:    res.Net,
	}, nil
}
