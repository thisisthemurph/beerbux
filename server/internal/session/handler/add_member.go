package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/sessionaccess"
	"beerbux/internal/common/useraccess"
	"beerbux/internal/session/command"
	"beerbux/pkg/send"
	"encoding/json"
	"errors"
	"fmt"
	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type AddSessionMemberHandler struct {
	userReader              useraccess.UserReader
	sessionReader           sessionaccess.SessionReader
	addSessionMemberCommand *command.AddSessionMemberCommand
	logger                  *slog.Logger
}

func NewAddSessionMemberHandler(
	userReader useraccess.UserReader,
	sessionReader sessionaccess.SessionReader,
	addSessionMemberCommand *command.AddSessionMemberCommand,
	logger *slog.Logger,
) *AddSessionMemberHandler {
	return &AddSessionMemberHandler{
		userReader:              userReader,
		sessionReader:           sessionReader,
		addSessionMemberCommand: addSessionMemberCommand,
		logger:                  logger,
	}
}

type AddMemberToSessionRequest struct {
	Username string `json:"username"`
}

// AddSessionMemberHandler godoc
// @Summary Add Member
// @Tags session
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 201 "Created"
// @Failure 400 {object} send.ErrorResponse "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /session/{sessionId}/member [post]
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

	currentMember, err := h.sessionReader.GetSessionMember(r.Context(), sessionID, c.Subject)
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

	if userIsMember, err := h.sessionReader.UserIsMemberOfSession(r.Context(), sessionID, userToAdd.ID); err != nil {
		send.InternalServerError(w, "There has been an issue determining if the user is already a member")
		return
	} else if userIsMember {
		// Return early if the user being added is already a member of the session.
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := h.addSessionMemberCommand.Execute(r.Context(), sessionID, userToAdd.ID, currentMember.ID); err != nil {
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
