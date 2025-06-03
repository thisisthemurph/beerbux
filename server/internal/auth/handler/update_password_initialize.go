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
	initializeUpdatePasswordCommand *command.InitializeUpdatePasswordCommand
	emailSender                     email.Sender
	logger                          *slog.Logger
}

func NewInitializeUpdatePasswordHandler(
	initializeUpdatePasswordCommand *command.InitializeUpdatePasswordCommand,
	emailSender email.Sender,
	logger *slog.Logger,
) *InitializePasswordUpdateHandler {
	return &InitializePasswordUpdateHandler{
		initializeUpdatePasswordCommand: initializeUpdatePasswordCommand,
		emailSender:                     emailSender,
		logger:                          logger,
	}
}

type InitializePasswordUpdateRequest struct {
	NewPassword string `json:"newPassword"`
}

// InitializePasswordUpdateHandler godoc
// @Summary Update password init
// @Description Initialize the update password process
// @Tags auth
// @Accept json
// @Produce json
// @Param login body InitializePasswordUpdateRequest true "New password request"
// @Success 200
// @Failure 400 {object} send.ErrorResponse "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
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

	result, err := h.initializeUpdatePasswordCommand.Execute(r.Context(), c.Subject, req.NewPassword)
	if err != nil {
		send.InternalServerError(w, "There has been an issue updating your password")
		return
	}

	h.sendUpdatePasswordEmail(c, result.OTP)
	w.WriteHeader(http.StatusOK)
}

func (h *InitializePasswordUpdateHandler) sendUpdatePasswordEmail(c claims.JWTClaims, otp string) {
	html, err := email.GenerateUpdatePasswordEmail(email.UpdatePasswordEmailData{
		Username:          c.Username,
		OTP:               otp,
		ExpirationMinutes: strconv.FormatInt(int64(command.OTPTimeToLiveMinutes), 10),
	})
	if err != nil {
		h.logger.Error("failed to generate update password email template", "error", err)
		return
	}
	if _, err := h.emailSender.Send(c.Email, "Password update request", html); err != nil {
		h.logger.Error("failed to send email", "error", err)
	}
}
