package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/common/useraccess"
	"beerbux/internal/user/command"
	"beerbux/pkg/send"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type UpdateUserHandler struct {
	updateUserCommand *command.UpdateUserCommand
	userReader        useraccess.UserReader
	logger            *slog.Logger
}

func NewUpdateUserHandler(
	updateUserCommand *command.UpdateUserCommand,
	userReader useraccess.UserReader,
	logger *slog.Logger,
) *UpdateUserHandler {
	return &UpdateUserHandler{
		updateUserCommand: updateUserCommand,
		userReader:        userReader,
		logger:            logger,
	}
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (h *UpdateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.BadRequest(w, "Failed to decode request")
		return
	}

	newUsername := strings.ToLower(req.Username)
	if strings.ToLower(c.Username) != newUsername {
		usernameTaken, err := h.userReader.UserWithUsernameExists(r.Context(), newUsername)
		if err != nil {
			send.InternalServerError(w, "Error checking if the username is already taken")
			return
		}
		if usernameTaken {
			send.BadRequest(w, fmt.Sprintf("Username %s is already taken", newUsername))
			return
		}
	}

	response, err := h.updateUserCommand.Execute(r.Context(), c.Subject, req.Name, newUsername)
	if err != nil {
		send.InternalServerError(w, "There has been an issue updating your details")
		return
	}

	send.JSON(w, response, http.StatusOK)
}
