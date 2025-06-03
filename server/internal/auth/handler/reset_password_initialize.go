package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/common/useraccess"
	"beerbux/pkg/email"
	"beerbux/pkg/send"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

type InitializePasswordResetHandler struct {
	initializePasswordResetCommand *command.InitializePasswordResetCommand
	userReader                     useraccess.UserReader
	emailSender                    email.Sender
	logger                         *slog.Logger
}

func NewInitializePasswordResetHandler(
	initializePasswordResetCommand *command.InitializePasswordResetCommand,
	userReader useraccess.UserReader,
	emailSender email.Sender,
	logger *slog.Logger,
) *InitializePasswordResetHandler {
	return &InitializePasswordResetHandler{
		initializePasswordResetCommand: initializePasswordResetCommand,
		userReader:                     userReader,
		emailSender:                    emailSender,
		logger:                         logger,
	}
}

type InitializePasswordResetRequest struct {
	Email string `json:"email"`
}

// InitializePasswordResetHandler godoc
// @Summary Password reset init
// @Description Initializes the password reset process.
// @Tags auth
// @Accept json
// @Produce json
// @Param login body InitializePasswordResetRequest true "User details"
// @Success 200
// @Failure 400 {object} send.ErrorResponse "Bad Request"
// @Failure 500 {object} send.ErrorResponse "Internal Server Error"
// @Router /auth/password/initialize-reset [post]
func (h *InitializePasswordResetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req InitializePasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.BadRequest(w, "Invalid request")
		return
	}

	user, err := h.userReader.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, useraccess.ErrUserNotFound) {
			// Return a status OK to indicate that an email has been sent if the email address could be found
			w.WriteHeader(http.StatusOK)
			return
		}
		h.logger.Error("failed to get user by email", "error", err)
		send.InternalServerError(w, "There has been an issue checking your email address, please try again")
		return
	}

	result, err := h.initializePasswordResetCommand.Execute(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to execute initialize password reset command", "error", err)
		send.InternalServerError(w, "There was an issue resetting your password, please try again")
		return
	}

	h.sendPasswordResetEmail(user.Email, user.Username, result.OTP)
	w.WriteHeader(http.StatusOK)
}

func (h *InitializePasswordResetHandler) sendPasswordResetEmail(emailAddress, username, otp string) {
	html, err := email.GenerateResetPasswordEmail(email.ResetPasswordEmailData{
		Username:          username,
		OTP:               otp,
		ExpirationMinutes: strconv.FormatInt(int64(command.OTPTimeToLiveMinutes), 10),
	})
	if err != nil {
		h.logger.Error("failed to generate password reset email template", "error", err)
		return
	}
	if _, err := h.emailSender.Send(emailAddress, "Password reset request", html); err != nil {
		h.logger.Error("failed to send email", "error", err)
	}
}
