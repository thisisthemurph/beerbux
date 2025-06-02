package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/sessionaccess"
	"beerbux/internal/session/command"
	"beerbux/internal/sse"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type CreateTransactionHandler struct {
	sessionReader            sessionaccess.SessionReader
	createTransactionCommand *command.CreateTransactionCommand
	logger                   *slog.Logger
	msgChan                  chan<- *sse.Message
}

func NewCreateTransactionHandler(
	sessionReader sessionaccess.SessionReader,
	createTransactionCommand *command.CreateTransactionCommand,
	logger *slog.Logger,
	msgChan chan<- *sse.Message,
) *CreateTransactionHandler {
	return &CreateTransactionHandler{
		sessionReader:            sessionReader,
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

	sessionID, ok := url.Path.GetUUID(r, "sessionId")
	if !ok {
		send.BadRequest(w, "Session ID is required")
		return
	}

	isMember, err := h.sessionReader.UserIsMemberOfSession(r.Context(), sessionID, c.Subject)
	if err != nil {
		send.InternalServerError(w, "There was an issue finding the session")
		return
	}
	if !isMember {
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

	if err := h.sendTransactionCreatedMessage(sessionID, c.Subject, createdTransaction.ID, transactionLines); err != nil {
		h.logger.Error("Failed to send session.transaction.created message", "error", err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CreateTransactionHandler) sendTransactionCreatedMessage(
	sessionID,
	creatorID,
	transactionID uuid.UUID,
	lines []command.TransactionLine,
) error {
	var total float64
	for _, line := range lines {
		total += line.Amount
	}

	data := TransactionCreatedMessage{
		TransactionID: transactionID,
		SessionID:     sessionID,
		CreatorID:     creatorID,
		Total:         total,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	h.msgChan <- sse.NewMessage("session.transaction.created", sessionID.String(), jsonData)
	return nil
}
