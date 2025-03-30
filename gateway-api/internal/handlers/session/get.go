package session

import (
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/dto"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type GetSessionByIdHandler struct {
	sessionClient sessionpb.SessionClient
}

func NewGetSessionByIdHandler(sessionClient sessionpb.SessionClient) *GetSessionByIdHandler {
	return &GetSessionByIdHandler{
		sessionClient: sessionClient,
	}
}

func (h *GetSessionByIdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionId, err := uuid.Parse(r.PathValue("sessionId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := h.sessionClient.GetSession(r.Context(), &sessionpb.GetSessionRequest{
		SessionId: sessionId.String(),
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch st.Code() {
		case codes.NotFound:
			send.Error(w, "Session not found", http.StatusBadRequest)
			return
		default:
			send.Error(w, "Could not fetch the session", http.StatusInternalServerError)
			return
		}
	}

	ssn := dto.SessionResponse{
		ID:           s.SessionId,
		Name:         s.Name,
		Total:        s.Total,
		IsActive:     s.IsActive,
		Members:      make([]dto.SessionMember, 0, len(s.Members)),
		Transactions: make([]dto.SessionTransaction, 0, len(s.Transactions)),
	}

	memberMap := make(map[string]dto.SessionMember, len(s.Members))
	for _, m := range s.Members {
		memberMap[m.UserId] = dto.SessionMember{
			ID:       m.UserId,
			Name:     m.Name,
			Username: m.Username,
		}

		ssn.Members = append(ssn.Members, dto.SessionMember{
			ID:       m.UserId,
			Name:     m.Name,
			Username: m.Username,
		})
	}

	for _, t := range s.Transactions {
		transactionMembers := make([]dto.SessionTransactionMember, 0, len(t.Lines))
		for _, l := range t.Lines {
			member, ok := memberMap[l.UserId]
			if !ok {
				// If the member is not found in the map, we will assign it generic unknown values.
				member = dto.SessionMember{
					Name:     "unknown",
					Username: "unknown",
				}
			}

			transactionMembers = append(transactionMembers, dto.SessionTransactionMember{
				ID:       l.UserId,
				Name:     member.Name,
				Username: member.Username,
				Amount:   l.Amount,
			})
		}

		ssn.Transactions = append(ssn.Transactions, dto.SessionTransaction{
			ID:        t.TransactionId,
			CreatorID: t.UserId,
			Total:     t.Total,
			Members:   transactionMembers,
		})
	}

	send.JSON(w, ssn, http.StatusOK)
}
