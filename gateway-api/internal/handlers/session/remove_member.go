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

type RemoveMemberFromSessionHandler struct {
	sessionClient sessionpb.SessionClient
}

func NewRemoveMemberFromSession(sessionClient sessionpb.SessionClient) http.Handler {
	return &RemoveMemberFromSessionHandler{
		sessionClient: sessionClient,
	}
}

func (h *RemoveMemberFromSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params, ok := h.getURLParams(w, r)
	if !ok {
		return
	}

	ssn, err := h.sessionClient.GetSession(r.Context(), &sessionpb.GetSessionRequest{
		SessionId: params.SessionID,
	})
	if err != nil {
		return
	}

	for _, m := range ssn.Members {
		if m.UserId == c.Subject {
			if !m.IsAdmin {
				send.Error(w, "You must be an admin to remove a member from a session", http.StatusUnauthorized)
				return
			}
		}
	}

	_, err = h.sessionClient.RemoveMemberFromSession(r.Context(), &sessionpb.RemoveMemberFromSessionRequest{
		SessionId: params.SessionID,
		UserId:    params.MemberID,
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

type RemoveMemberFromSessionURLParams struct {
	SessionID string
	MemberID  string
}

func (h *RemoveMemberFromSessionHandler) getURLParams(w http.ResponseWriter, r *http.Request) (RemoveMemberFromSessionURLParams, bool) {
	sessionID, ok := url.GetIDFromPath(w, r, "sessionId")
	if !ok {
		return RemoveMemberFromSessionURLParams{}, false
	}

	memberID, ok := url.GetIDFromPath(w, r, "memberId")
	if !ok {
		return RemoveMemberFromSessionURLParams{}, false
	}

	return RemoveMemberFromSessionURLParams{
		SessionID: sessionID,
		MemberID:  memberID,
	}, true
}
