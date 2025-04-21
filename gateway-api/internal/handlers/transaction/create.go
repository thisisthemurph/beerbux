package transaction

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/transaction-service/protos/transactionpb"
)

type CreateTransactionHandler struct {
	logger            *slog.Logger
	transactionClient transactionpb.TransactionClient
}

func NewCreateTransactionHandler(logger *slog.Logger, transactionClient transactionpb.TransactionClient) *CreateTransactionHandler {
	return &CreateTransactionHandler{
		logger:            logger,
		transactionClient: transactionClient,
	}
}

func (h *CreateTransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID := r.PathValue("sessionId")
	if _, err := uuid.Parse(sessionID); err != nil {
		send.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var transactions map[string]int
	if err := json.NewDecoder(r.Body).Decode(&transactions); err != nil {
		send.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	memberAmounts := make([]*transactionpb.MemberAmount, 0, len(transactions))
	for userID, amount := range transactions {
		if userID == c.Subject || amount <= 0 {
			continue
		}

		memberAmounts = append(memberAmounts, &transactionpb.MemberAmount{
			UserId: userID,
			Amount: float64(amount),
		})
	}

	_, err := h.transactionClient.CreateTransaction(r.Context(), &transactionpb.CreateTransactionRequest{
		CreatorId:     c.Subject,
		SessionId:     sessionID,
		MemberAmounts: memberAmounts,
	})

	if err != nil {
		h.logger.Error("Failed to create transaction", "error", err)
		send.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
