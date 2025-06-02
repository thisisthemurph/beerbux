package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/sessionaccess"
	"beerbux/internal/session/command"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type UpdateSessionActiveStateHandler struct {
	sessionReader                   sessionaccess.SessionReader
	updateSessionActiveStateCommand *command.UpdateSessionActiveStateCommand
	logger                          *slog.Logger
}

func NewUpdateSessionActiveStateHandler(
	sr sessionaccess.SessionReader,
	updateSessionActiveStateCommand *command.UpdateSessionActiveStateCommand,
	logger *slog.Logger,
) *UpdateSessionActiveStateHandler {
	return &UpdateSessionActiveStateHandler{
		sessionReader:                   sr,
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

	currentMember, err := h.sessionReader.GetSessionMember(r.Context(), params.SessionID, c.Subject)
	if errors.Is(err, sessionaccess.ErrMemberNotFound) {
		send.Unauthorized(w, "You are not a member of this session")
		return
	} else if err != nil {
		send.InternalServerError(w, "There has been an issue determining if you are a member of the session.")
		return
	} else if !currentMember.IsAdmin {
		send.Unauthorized(w, "You must be an admin to add a member to a session")
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
