package session

import (
	"fmt"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/url"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"log/slog"
	"net/http"
)

type UpdateSessionActiveStateHandler struct {
	logger        *slog.Logger
	sessionClient sessionpb.SessionClient
}

func NewUpdateSessionActiveStateHandler(logger *slog.Logger, sessionClient sessionpb.SessionClient) http.Handler {
	return &UpdateSessionActiveStateHandler{
		logger:        logger,
		sessionClient: sessionClient,
	}
}

func (h *UpdateSessionActiveStateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params, ok := h.getURLParams(w, r)
	if !ok {
		return
	}

	isAdminMember, err := memberIsAdmin(r.Context(), h.sessionClient, params.SessionID, c.Subject)
	if !isAdminMember || err != nil {
		send.Error(w, fmt.Sprintf("You do not have permission to %s this session", params.Command), http.StatusUnauthorized)
		return
	}

	_, err = h.sessionClient.UpdateSessionActiveState(r.Context(), &sessionpb.UpdateSessionActiveStateRequest{
		SessionId:     params.SessionID,
		IsActive:      params.NewIsActiveState,
		PerformedById: c.Subject,
	})
	if err != nil {
		h.logger.Error("Error updating session active state", "error", err)
		send.Error(w, fmt.Sprintf("Failed to %s the session", params.Command), http.StatusInternalServerError)
		return
	}
}

type UpdateSessionActiveStateURLParams struct {
	SessionID        string
	Command          string
	NewIsActiveState bool
}

func (h *UpdateSessionActiveStateHandler) getURLParams(w http.ResponseWriter, r *http.Request) (UpdateSessionActiveStateURLParams, bool) {
	sessionID, ok := url.GetIDFromPath(w, r, "sessionId")
	if !ok {
		return UpdateSessionActiveStateURLParams{}, false
	}

	command, ok := url.GetTextFromPath(w, r, "command")
	if !ok {
		return UpdateSessionActiveStateURLParams{}, false
	}

	if command != "activate" && command != "deactivate" {
		h.logger.Warn("Unexpected command, expected activate or deactivate")
		return UpdateSessionActiveStateURLParams{}, false
	}

	return UpdateSessionActiveStateURLParams{
		SessionID:        sessionID,
		Command:          command,
		NewIsActiveState: command == "activate",
	}, true
}
