package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/session/command"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/query"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type UpdateSessionActiveStateHandler struct {
	getSessionQuery                 *query.GetSessionQuery
	updateSessionActiveStateCommand *command.UpdateSessionActiveStateCommand
	logger                          *slog.Logger
}

func NewUpdateSessionActiveStateHandler(
	getSessionQuery *query.GetSessionQuery,
	updateSessionActiveStateCommand *command.UpdateSessionActiveStateCommand,
	logger *slog.Logger,
) *UpdateSessionActiveStateHandler {
	return &UpdateSessionActiveStateHandler{
		getSessionQuery:                 getSessionQuery,
		updateSessionActiveStateCommand: updateSessionActiveStateCommand,
		logger:                          logger,
	}
}

type UpdateSessionActiveStateURLParams struct {
	SessionID        uuid.UUID
	Command          string
	NewIsActiveState bool
}

func (h *UpdateSessionActiveStateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params, ok := h.getURLParams(r)
	if !ok {
		send.BadRequest(w, "Missing URL parameters")
		return
	}

	s, err := h.getSessionQuery.Execute(r.Context(), params.SessionID)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "The session could not be found")
			return
		}
		send.InternalServerError(w, "There was an issue finding the session")
		return
	}

	if !s.IsMember(c.Subject) {
		send.Unauthorized(w, "You are not a member of this session")
		return
	}
	if !s.IsAdminMember(c.Subject) {
		send.Unauthorized(w, fmt.Sprintf("You must be an admin to %s this session", params.Command))
		return
	}

	err = h.updateSessionActiveStateCommand.Execute(r.Context(), params.SessionID, c.Subject, params.NewIsActiveState)
	if err != nil {
		h.logger.Error("failed to update the session active state", "error", err)
		send.InternalServerError(w, "There has been an issue updating the session active state")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UpdateSessionActiveStateHandler) getURLParams(r *http.Request) (UpdateSessionActiveStateURLParams, bool) {
	sessionID, ok := url.Path.GetUUID(r, "sessionId")
	if !ok {
		return UpdateSessionActiveStateURLParams{}, false
	}

	cmd, ok := url.Path.GetString(r, "command")
	if !ok {
		return UpdateSessionActiveStateURLParams{}, false
	}

	if cmd != "activate" && cmd != "deactivate" {
		h.logger.Warn("Unexpected command, expected activate or deactivate")
		return UpdateSessionActiveStateURLParams{}, false
	}

	return UpdateSessionActiveStateURLParams{
		SessionID:        sessionID,
		Command:          cmd,
		NewIsActiveState: cmd == "activate",
	}, true
}
