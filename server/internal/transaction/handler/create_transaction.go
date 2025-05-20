package handler

import (
	"beerbux/internal/common/claims"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/query"
	"beerbux/internal/sse"
	"beerbux/internal/transaction/command"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type CreateTransactionHandler struct {
	getSessionQuery          *query.GetSessionQuery
	createTransactionCommand *command.CreateTransactionCommand
	logger                   *slog.Logger
	msgChan                  chan<- *sse.Message
}

func NewCreateTransactionHandler(
	getSessionQuery *query.GetSessionQuery,
	createTransactionCommand *command.CreateTransactionCommand,
	logger *slog.Logger,
	msgChan chan<- *sse.Message,
) *CreateTransactionHandler {
	return &CreateTransactionHandler{
		getSessionQuery:          getSessionQuery,
		createTransactionCommand: createTransactionCommand,
		logger:                   logger,
		msgChan:                  msgChan,
	}
}

type TransactionCreatedMessage struct {
	TransactionID uuid.UUID `json:"transactionId"`
	SessionID     uuid.UUID `json:"sessionID"`
	CreatorID     uuid.UUID `json:"creatorId"`
	Total         float64   `json:"total"`
}

func (h *CreateTransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, ok := url.GetUUIDFromPath(r, "sessionId")
	if !ok {
		send.BadRequest(w, "Session ID is required")
		return
	}

	session, err := h.getSessionQuery.Execute(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "The session could not be found")
			return
		}
		send.InternalServerError(w, "There was an issue finding the session")
		return
	}
	if !session.IsMember(c.Subject) {
		send.Unauthorized(w, "You are not a member of the session")
		return
	}

	var transactionLineRecords map[uuid.UUID]float64
	if err := json.NewDecoder(r.Body).Decode(&transactionLineRecords); err != nil {
		send.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	transactionLines := make([]command.TransactionLine, 0, len(transactionLineRecords))
	for memberID, amount := range transactionLineRecords {
		if memberID == c.Subject || amount <= 0 {
			continue
		}

		transactionLines = append(transactionLines, command.TransactionLine{
			MemberID: memberID,
			Amount:   amount,
		})
	}

	createdTransaction, err := h.createTransactionCommand.Execute(r.Context(), command.CreateTransactionRequest{
		SessionID: sessionID,
		CreatorID: c.Subject,
		Lines:     transactionLines,
	})
	if err != nil {
		h.logger.Error("failed to create transaction", "error", err)
		send.InternalServerError(w, "There was an issue creating the transaction")
		return
	}

	if err := h.sendTransactionCreatedMessage(createdTransaction.ID, sessionID, c.Subject); err != nil {
		h.logger.Error("Failed to send session.transaction.created message", "error", err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CreateTransactionHandler) sendTransactionCreatedMessage(
	transactionID,
	sessionID,
	creatorID uuid.UUID,
) error {
	data := TransactionCreatedMessage{
		TransactionID: transactionID,
		SessionID:     sessionID,
		CreatorID:     creatorID,
		Total:         0,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	h.msgChan <- sse.NewMessage("session.transaction.created", sessionID.String(), jsonData)
	return nil
}
