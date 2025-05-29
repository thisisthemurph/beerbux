package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/common/claims"
	"beerbux/pkg/email"
	"beerbux/pkg/send"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type InitializePasswordUpdateHandler struct {
	initializePasswordUpdateCommand *command.InitializePasswordResetCommand
	emailSender                     email.Sender
	logger                          *slog.Logger
}

func NewInitializePasswordUpdateHandler(
	initializePasswordResetCommand *command.InitializePasswordResetCommand,
	emailSender email.Sender,
	logger *slog.Logger,
) *InitializePasswordUpdateHandler {
	return &InitializePasswordUpdateHandler{
		initializePasswordUpdateCommand: initializePasswordResetCommand,
		emailSender:                     emailSender,
		logger:                          logger,
	}
}

type InitializePasswordUpdateRequest struct {
	NewPassword string `json:"newPassword"`
}

func (h *InitializePasswordUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req InitializePasswordUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.BadRequest(w, "Failed to read request body")
		return
	}

	result, err := h.initializePasswordUpdateCommand.Execute(r.Context(), c.Subject, req.NewPassword)
	if err != nil {
		send.InternalServerError(w, "There has been an issue updating your password")
		return
	}

	h.sendPasswordResetEmail(c, result.OTP)
	w.WriteHeader(http.StatusOK)
}

func (h *InitializePasswordUpdateHandler) sendPasswordResetEmail(c claims.JWTClaims, otp string) {
	html, err := email.GeneratePasswordResetEmail(email.PasswordResetEmailData{
		Username:          c.Username,
		OTP:               otp,
		ExpirationMinutes: strconv.FormatInt(int64(command.OTPTimeToLiveMinutes), 10),
	})
	if err != nil {
		h.logger.Error("failed to generate password reset email template", "error", err)
		return
	}
	if _, err := h.emailSender.Send(c.Email, "Password reset request", html); err != nil {
		h.logger.Error("failed to send email", "error", err)
	}
}
