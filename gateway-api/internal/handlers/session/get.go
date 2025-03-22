package session

import (
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
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

	ssn := SessionResponse{
		ID:       s.SessionId,
		Name:     s.Name,
		IsActive: s.IsActive,
		Members:  make([]SessionMember, 0, len(s.Members)),
	}

	for _, m := range s.Members {
		ssn.Members = append(ssn.Members, SessionMember{
			ID:       m.UserId,
			Name:     m.Name,
			Username: m.Username,
		})
	}

	send.JSON(w, ssn, http.StatusOK)
}
