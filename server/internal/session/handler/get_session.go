package handler

import (
	"beerbux/internal/common/claims"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/query"
	"beerbux/pkg/send"
	"errors"
	"github.com/google/uuid"
	"github.com/thisisthemurph/fn"
	"log/slog"
	"net/http"
)

type GetSessionHandler struct {
	getSessionQuery *query.GetSessionQuery
	logger          *slog.Logger
}

func NewGetSessionHandler(getSessionQuery *query.GetSessionQuery, logger *slog.Logger) *GetSessionHandler {
	return &GetSessionHandler{
		getSessionQuery: getSessionQuery,
		logger:          logger,
	}
}

type GetSessionResponse struct {
	ID           uuid.UUID                  `json:"id"`
	Name         string                     `json:"name"`
	Total        float64                    `json:"total"`
	IsActive     bool                       `json:"isActive"`
	Members      []GetSessionResponseMember `json:"members"`
	Transactions []query.SessionTransaction `json:"transactions"`
}

type GetSessionResponseMember struct {
	query.SessionMember
	TransactionSummary TransactionSummary `json:"transactionSummary"`
}

type TransactionSummary struct {
	Credit float64 `json:"credit"`
	Debit  float64 `json:"debit"`
}

func (h *GetSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, err := uuid.Parse(r.PathValue("sessionId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := h.getSessionQuery.Execute(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "Session not found")
			return
		}

		h.logger.Error("failed to fetch the session", "session", sessionID, "error", err)
		send.InternalServerError(w, "There was an issie fetching the session")
		return
	}

	if err := h.validateUserAgainstSessionMembers(c.Subject, s.Members); err != nil {
		send.Unauthorized(w, err.Error())
		return
	}

	send.JSON(w, h.buildResponse(s), http.StatusOK)
}

func (h *GetSessionHandler) buildResponse(s *query.SessionResponse) GetSessionResponse {
	members := fn.Map(s.Members, func(m query.SessionMember) GetSessionResponseMember {
		return GetSessionResponseMember{
			SessionMember:      m,
			TransactionSummary: h.calculateTransactionSummaryForMember(s.Transactions, m.ID),
		}
	})

	return GetSessionResponse{
		ID:           s.ID,
		Name:         s.Name,
		Total:        s.Total,
		IsActive:     s.IsActive,
		Members:      members,
		Transactions: s.Transactions,
	}
}

func (h *GetSessionHandler) calculateTransactionSummaryForMember(transactions []query.SessionTransaction, memberID uuid.UUID) TransactionSummary {
	summary := TransactionSummary{}
	for _, t := range transactions {
		// This member created the transaction.
		if t.UserID == memberID {
			summary.Debit += t.Total
			continue
		}

		for _, line := range t.Lines {
			// This member is a participant in the transaction.
			if line.UserID == memberID {
				summary.Credit += line.Amount
				continue
			}
		}
	}

	return summary
}

func (h *GetSessionHandler) validateUserAgainstSessionMembers(userID uuid.UUID, members []query.SessionMember) error {
	userIsMember := false
	for _, m := range members {
		if m.ID == userID {
			userIsMember = true
			if m.IsDeleted {
				return errors.New("you were removed from this session and do not have permission to access it")
			}
			break
		}
	}
	if !userIsMember {
		return errors.New("you are not a member of this session")
	}
	return nil
}
