package handler

import (
	"beerbux/internal/common/claims"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/query"
	"beerbux/pkg/send"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type GetSessionHandler struct {
	getSessionQuery *query.GetSessionQuery
	logger          *slog.Logger
}

func NewGetSessionHandler(getSessionQuery *query.GetSessionQuery, logger *slog.Logger) *GetSessionHandler {
	return &GetSessionHandler{
		getSessionQuery: getSessionQuery,
		logger:          logger,
	}
}

func (h *GetSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, err := uuid.Parse(r.PathValue("sessionId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := h.getSessionQuery.Execute(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "Session not found")
			return
		}

		h.logger.Error("failed to fetch the session", "session", sessionID, "error", err)
		send.InternalServerError(w, "There was an issie fetching the session")
		return
	}

	if err := h.validateUserAgainstSessionMembers(c.Subject, s.Members); err != nil {
		send.Unauthorized(w, err.Error())
		return
	}

	send.JSON(w, s, http.StatusOK)
}

func (h *GetSessionHandler) validateUserAgainstSessionMembers(userID uuid.UUID, members []query.SessionMember) error {
	userIsMember := false
	for _, m := range members {
		if m.ID == userID {
			userIsMember = true
			if m.IsDeleted {
				return errors.New("you were removed from this session and do not have permission to access it")
			}
			break
		}
	}
	if !userIsMember {
		return errors.New("you are not a member of this session")
	}
	return nil
}
