package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/sessionaccess"
	"beerbux/internal/session/command"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type RemoveSessionMemberHandler struct {
	sessionReader              sessionaccess.SessionReader
	removeSessionMemberCommand *command.RemoveSessionMemberCommand
	logger                     *slog.Logger
}

func NewRemoveSessionMemberHandler(
	sr sessionaccess.SessionReader,
	removeSessionMemberCommand *command.RemoveSessionMemberCommand,
	logger *slog.Logger,
) *RemoveSessionMemberHandler {
	return &RemoveSessionMemberHandler{
		sessionReader:              sr,
		removeSessionMemberCommand: removeSessionMemberCommand,
		logger:                     logger,
	}
}

func (h *RemoveSessionMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params, ok := h.getURLParams(r)
	if !ok {
		send.BadRequest(w, "Session and member IDs are required")
		return
	}

	currentMember, err := h.sessionReader.GetSessionMember(r.Context(), params.SessionID, c.Subject)
	if err != nil {
		if errors.Is(err, sessionaccess.ErrMemberNotFound) {
			send.Unauthorized(w, "You are not a member of this session")
			return
		}
		send.InternalServerError(w, "There was an issue finding your session")
		return
	}
	if !currentMember.IsAdmin {
		send.Unauthorized(w, "You must be an admin to remove a member from the session")
		return
	}

	if err := h.removeSessionMemberCommand.Execute(r.Context(), params.SessionID, params.MemberID, c.Subject); err != nil {
		if errors.Is(err, sessionErr.ErrSessionMemberNotFound) {
			send.NotFound(w, "The session could not be found")
		} else if errors.Is(err, sessionErr.ErrSessionMustHaveAtLeastOneMember) {
			send.BadRequest(w, "Could not remove the member, the session must have at least one member")
		} else if errors.Is(err, sessionErr.ErrSessionMustHaveAtLeastOneAdmin) {
			send.BadRequest(w, "Could not remove the member, the session must have at least one admin")
		} else {
			h.logger.Error("failed to remove a member from the session", "memberToRemove", params.MemberID, "session", params.SessionID, "actor", c.Subject, "error", err)
			send.InternalServerError(w, "There has been an issue removing the member from the session")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

type RemoveSessionMemberURLParams struct {
	SessionID uuid.UUID
	MemberID  uuid.UUID
}

func (h *RemoveSessionMemberHandler) getURLParams(r *http.Request) (RemoveSessionMemberURLParams, bool) {
	sessionID, ok := url.Path.GetUUID(r, "sessionId")
	if !ok {
		return RemoveSessionMemberURLParams{}, false
	}

	memberID, ok := url.Path.GetUUID(r, "memberId")
	if !ok {
		return RemoveSessionMemberURLParams{}, false
	}

	return RemoveSessionMemberURLParams{
		SessionID: sessionID,
		MemberID:  memberID,
	}, true
}
