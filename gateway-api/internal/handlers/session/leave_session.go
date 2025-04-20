package session

import (
	"net/http"

	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/url"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LeaveSessionHandler struct {
	sessionClient sessionpb.SessionClient
}

func NewLeaveSessionHandler(sessionClient sessionpb.SessionClient) http.Handler {
	return &LeaveSessionHandler{
		sessionClient: sessionClient,
	}
}

func (h *LeaveSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params, ok := h.getURLParams(w, r)
	if !ok {
		return
	}

	_, err := h.sessionClient.RemoveMemberFromSession(r.Context(), &sessionpb.RemoveMemberFromSessionRequest{
		SessionId: params.SessionID,
		UserId:    c.Subject,
	})
	if err != nil {
		sc, ok := status.FromError(err)
		if ok {
			switch sc.Code() {
			case codes.FailedPrecondition:
				send.Error(w, sc.Message(), http.StatusBadRequest)
			case codes.Internal:
				send.Error(w, "Failed to remove user from session", http.StatusInternalServerError)
			}
			return
		}

		send.Error(w, "Unexpected error removing member from session", http.StatusInternalServerError)
		return
	}
}

type LeaveSessionHandlerURLParams struct {
	SessionID string
}

func (h *LeaveSessionHandler) getURLParams(w http.ResponseWriter, r *http.Request) (LeaveSessionHandlerURLParams, bool) {
	sessionID, ok := url.GetIDFromPath(w, r, "sessionId")
	if !ok {
		return LeaveSessionHandlerURLParams{}, false
	}

	return LeaveSessionHandlerURLParams{
		SessionID: sessionID,
	}, true
}
