package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/session/command"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"errors"
	"github.com/google/uuid"
	"net/http"
)

type LeaveSessionHandler struct {
	removeSessionMemberCommand *command.RemoveSessionMemberCommand
}

func (h *LeaveSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, ok := url.GetUUIDFromPath(r, "sessionId")
	if !ok {
		send.BadRequest(w, "Session ID is required")
		return
	}

	err := h.removeSessionMemberCommand.Execute(r.Context(), sessionID, c.Subject, c.Subject)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "Session not found")
		} else if errors.Is(err, sessionErr.ErrSessionMemberNotFound) {
			send.NotFound(w, "Session member not found")
		} else if errors.Is(err, sessionErr.ErrSessionMustHaveAtLeastOneMember) {
			send.BadRequest(w, "You cannot leave the session if you are the only member")
		} else if errors.Is(err, sessionErr.ErrSessionMustHaveAtLeastOneAdmin) {
			send.BadRequest(w, "You cannot leave the session if you are the only admin member")
		} else {
			send.InternalServerError(w, "There was an issue leaving the session")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

type LeaveSessionHandlerURLParams struct {
	SessionID uuid.UUID
}

func (h *LeaveSessionHandler) getURLParams(r *http.Request) (LeaveSessionHandlerURLParams, bool) {
	sessionID, ok := url.GetUUIDFromPath(r, "sessionId")
	if !ok {
		return LeaveSessionHandlerURLParams{}, false
	}

	return LeaveSessionHandlerURLParams{
		SessionID: sessionID,
	}, true
}
