package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/pkg/send"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type ResetPasswordHandler struct {
	resetPasswordCommand *command.ResetPasswordCommand
	logger               *slog.Logger
}

func NewResetPasswordHandler(resetPasswordCommand *command.ResetPasswordCommand, logger *slog.Logger) *ResetPasswordHandler {
	return &ResetPasswordHandler{
		resetPasswordCommand: resetPasswordCommand,
		logger:               logger,
	}
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
	OTP         string `json:"otp"`
}

// NewResetPasswordHandler godoc
// @Summary Reset Password
// @Tags auth
// @Accept json
// @Produce json
// @Param login body ResetPasswordRequest true "Password reset request"
// @Failure 400 {object} send.ErrorResponse "Bad Request"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /auth/password/reset [put]
func (h *ResetPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.BadRequest(w, "Failed to decode request body")
		return
	}

	if passwordErr := validatePassword(req.NewPassword); passwordErr != nil {
		send.BadRequest(w, fmt.Sprintf("Invalid password: %s", passwordErr))
	}

	if err := h.resetPasswordCommand.Execute(r.Context(), req.Email, req.OTP, req.NewPassword); err != nil {
		h.handleResetPasswordError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ResetPasswordHandler) handleResetPasswordError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, command.ErrPasswordResetNotInitialized):
		send.BadRequest(w, "A password reset was not requested for this account")
	case errors.Is(err, command.ErrIncorrectOTP):
		send.BadRequest(w, "The provided OTP is incorrect")
	default:
		h.logger.Error("failed to reset the password", "error", err)
		send.InternalServerError(w, "There has been an issue resetting the password")
	}
}
