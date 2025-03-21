package session

import (
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"net/http"
)

type ListSessionsForUserHandler struct {
	sessionClient sessionpb.SessionClient
}

func NewListSessionsForUserHandler(sessionClient sessionpb.SessionClient) *ListSessionsForUserHandler {
	return &ListSessionsForUserHandler{
		sessionClient: sessionClient,
	}
}

type ListSessionsForUserResponse struct {
}

func (h *ListSessionsForUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}
