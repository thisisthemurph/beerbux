package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/session/command"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/query"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"errors"
	"github.com/google/uuid"
	"github.com/thisisthemurph/fn"
	"log/slog"
	"net/http"
)

type RemoveSessionMemberHandler struct {
	getSessionQuery            *query.GetSessionQuery
	removeSessionMemberCommand *command.RemoveSessionMemberCommand
	logger                     *slog.Logger
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

	s, err := h.getSessionQuery.Execute(r.Context(), params.SessionID)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "The session could not be found")
			return
		}
		send.InternalServerError(w, "There was an issue finding your session")
		return
	}

	currentMember, exists := fn.First(s.Members, func(member query.SessionMember) bool {
		return member.ID == c.Subject
	})
	if !exists {
		send.Unauthorized(w, "You are not a member of this session")
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
	sessionID, ok := url.GetUUIDFromPath(r, "sessionId")
	if !ok {
		return RemoveSessionMemberURLParams{}, false
	}

	memberID, ok := url.GetUUIDFromPath(r, "memberId")
	if !ok {
		return RemoveSessionMemberURLParams{}, false
	}

	return RemoveSessionMemberURLParams{
		SessionID: sessionID,
		MemberID:  memberID,
	}, true
}
