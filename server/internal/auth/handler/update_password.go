package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/common/claims"
	"beerbux/pkg/send"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type UpdatePasswordHandler struct {
	updatePasswordCommand *command.UpdatePasswordCommand
	logger                *slog.Logger
}

func NewUpdatePasswordHandler(updatePasswordCommand *command.UpdatePasswordCommand, logger *slog.Logger) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{
		updatePasswordCommand: updatePasswordCommand,
		logger:                logger,
	}
}

type UpdatePasswordRequest struct {
	OTP string `json:"otp"`
}

func (h *UpdatePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.BadRequest(w, "Failed to decode the request body")
		return
	}

	if err := h.updatePasswordCommand.Execute(r.Context(), c.Subject, req.OTP); err != nil {
		h.handleUpdatePasswordError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UpdatePasswordHandler) handleUpdatePasswordError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, command.ErrUserNotFound):
		send.NotFound(w, "User not found")
	case errors.Is(err, command.ErrIncorrectOTP):
		send.BadRequest(w, "The provided OTP is incorrect")
	case errors.Is(err, command.ErrOTPExpired):
		send.BadRequest(w, "Your OTP has expired, please start the process again")
	default:
		send.InternalServerError(w, "There was an issue updating your password")
	}
}
