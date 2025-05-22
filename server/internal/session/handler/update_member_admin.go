package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/session/command"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/query"
	"beerbux/pkg/send"
	"beerbux/pkg/url"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type UpdateSessionMemberAdminHandler struct {
	getSessionQuery               *query.GetSessionQuery
	updateMemberAdminStateCommand *command.UpdateSessionMemberAdminStateCommand
	logger                        *slog.Logger
}

func NewUpdateSessionMemberAdminHandler(
	getSessionQuery *query.GetSessionQuery,
	updateMemberAdminStateCommand *command.UpdateSessionMemberAdminStateCommand,
	logger *slog.Logger,
) *UpdateSessionMemberAdminHandler {
	return &UpdateSessionMemberAdminHandler{
		getSessionQuery:               getSessionQuery,
		updateMemberAdminStateCommand: updateMemberAdminStateCommand,
		logger:                        logger,
	}
}

type UpdateSessionMemberAdminRequest struct {
	NewAdminState bool `json:"newAdminState"`
}

func (h *UpdateSessionMemberAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params, ok := h.getURLParams(r)
	if !ok {
		send.BadRequest(w, "Missing params")
		return
	}

	var req UpdateSessionMemberAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("error decoding body", "error", err)
		send.BadRequest(w, "Failed to decode request")
		return
	}

	if params.MemberID == c.Subject {
		send.BadRequest(w, "You cannot update your own admin status")
		return
	}

	session, err := h.getSessionQuery.Execute(r.Context(), params.SessionID)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "The session could not be found")
			return
		}
		send.InternalServerError(w, "There was an issue finding the session")
		return
	}

	if !session.IsMember(c.Subject) {
		send.Unauthorized(w, "You are not a member of this session")
		return
	}
	if !session.IsAdminMember(c.Subject) {
		send.Unauthorized(w, "You must be an admin to add a member to the session")
		return
	}
	if !session.IsMember(params.MemberID) {
		send.BadRequest(w, "Member to update not found")
		return
	}

	err = h.updateMemberAdminStateCommand.Execute(r.Context(), session.ID, params.MemberID, req.NewAdminState, c.Subject)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionMustHaveAtLeastOneAdmin) {
			send.BadRequest(w, "A session must have at least one admin member")
			return
		}
		send.InternalServerError(w, "There has been an issue updating the admin status")
		return
	}

	w.WriteHeader(http.StatusOK)
}

type UpdateSessionMemberAdminURLParams struct {
	SessionID uuid.UUID
	MemberID  uuid.UUID
}

func (h *UpdateSessionMemberAdminHandler) getURLParams(r *http.Request) (UpdateSessionMemberAdminURLParams, bool) {
	sessionID, ok := url.Path.GetUUID(r, "sessionId")
	if !ok {
		return UpdateSessionMemberAdminURLParams{}, false
	}
	memberID, ok := url.Path.GetUUID(r, "memberId")
	if !ok {
		return UpdateSessionMemberAdminURLParams{}, false
	}

	return UpdateSessionMemberAdminURLParams{
		SessionID: sessionID,
		MemberID:  memberID,
	}, true
}
