package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/useraccess"
	"beerbux/internal/session/command"
	sessionErr "beerbux/internal/session/errors"
	"beerbux/internal/session/query"
	"beerbux/pkg/send"
	"encoding/json"
	"errors"
	"fmt"
	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/thisisthemurph/fn"
	"log/slog"
	"net/http"
)

type AddSessionMemberHandler struct {
	getSessionQuery         *query.GetSessionQuery
	userReader              useraccess.UserReader
	addSessionMemberCommand *command.AddSessionMemberCommand
	logger                  *slog.Logger
}

func NewAddSessionMemberHandler(
	getSessionQuery *query.GetSessionQuery,
	userReader useraccess.UserReader,
	addSessionMemberCommand *command.AddSessionMemberCommand,
	logger *slog.Logger,
) *AddSessionMemberHandler {
	return &AddSessionMemberHandler{
		getSessionQuery:         getSessionQuery,
		userReader:              userReader,
		addSessionMemberCommand: addSessionMemberCommand,
		logger:                  logger,
	}
}

type AddMemberToSessionRequest struct {
	Username string `json:"username"`
}

func (h *AddSessionMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, err := uuid.Parse(r.PathValue("sessionId"))
	if err != nil {
		send.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var req AddMemberToSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	s, err := h.getSessionQuery.Execute(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sessionErr.ErrSessionNotFound) {
			send.NotFound(w, "The session could not be found")
			return
		}
		h.logger.Error("failed to find session when adding member", "member", req.Username, "session", sessionID, "error", err)
		send.InternalServerError(w, "There was an issue finding the session")
		return
	}

	currentMember, exists := fn.First(s.Members, func(m query.SessionMember) bool {
		return m.ID == c.Subject
	})
	if !exists {
		send.Unauthorized(w, "You are not a member of this session")
		return
	}
	if !currentMember.IsAdmin {
		send.Unauthorized(w, "You must be an admin to add a member to a session")
		return
	}

	userToAdd, err := h.userReader.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		if errors.Is(err, useraccess.ErrUserNotFound) {
			send.NotFound(w, fmt.Sprintf("User %s not found", req.Username))
			return
		}
		h.logger.Error("failed to find the user to add to the session", "username", req.Username, "session", sessionID, "error", err)
		send.InternalServerError(w, "There has been an issue finding the user to add")
		return
	}

	// Return early if the user being added is already a member of the session.
	for _, m := range s.Members {
		if m.ID == userToAdd.ID {
			if !m.IsDeleted {
				w.WriteHeader(http.StatusOK)
				return
			}
			break
		}
	}

	if err := h.addSessionMemberCommand.Execute(r.Context(), s.ID, userToAdd.ID, currentMember.ID); err != nil {
		h.logger.Error("failed to add the user to the session", "error", err)
		send.InternalServerError(w, fmt.Sprintf("There has been an issue adding %s to the session", userToAdd.Username))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (r AddMemberToSessionRequest) Validate() error {
	return oz.ValidateStruct(&r,
		oz.Field(&r.Username,
			oz.Required.Error("The member's username is required"),
		),
	)
}
