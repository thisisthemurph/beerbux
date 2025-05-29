package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/user/command"
	"beerbux/pkg/send"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type UpdateUserHandler struct {
	updateUserCommand *command.UpdateUserCommand
	logger            *slog.Logger
}

func NewUpdateUserHandler(
	updateUserCommand *command.UpdateUserCommand,
	logger *slog.Logger,
) *UpdateUserHandler {
	return &UpdateUserHandler{
		updateUserCommand: updateUserCommand,
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

	username := strings.ToLower(req.Username)
	response, err := h.updateUserCommand.Execute(r.Context(), c.Subject, req.Name, username)
	if err != nil {
		if errors.Is(err, command.ErrUsernameExists) {
			msg := fmt.Sprintf("The username %s is already taken", username)
			send.BadRequest(w, msg)
			return
		}
		send.InternalServerError(w, "There has been an issue updating your details")
		return
	}

	send.JSON(w, response, http.StatusOK)
}
