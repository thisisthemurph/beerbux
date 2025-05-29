package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/common/claims"
	"beerbux/internal/common/useraccess"
	"beerbux/pkg/email"
	"beerbux/pkg/send"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type InitializeEmailUpdateHandler struct {
	initializeUpdateEmailCommand *command.InitializeUpdateEmailCommand
	userReader                   useraccess.UserReader
	emailSender                  email.Sender
	logger                       *slog.Logger
}

func NewInitializeEmailUpdateHandler(
	initializeUpdateEmailCommand *command.InitializeUpdateEmailCommand,
	userReader useraccess.UserReader,
	emailSender email.Sender,
	logger *slog.Logger,
) *InitializeEmailUpdateHandler {
	return &InitializeEmailUpdateHandler{
		initializeUpdateEmailCommand: initializeUpdateEmailCommand,
		userReader:                   userReader,
		emailSender:                  emailSender,
		logger:                       logger,
	}
}

type InitializeEmailUpdateRequest struct {
	NewEmail string `json:"newEmail"`
}

func (h *InitializeEmailUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req InitializeEmailUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.BadRequest(w, "Failed to read request body")
		return
	}

	if req.NewEmail == c.Email {
		send.BadRequest(w, "This is your current email address")
		return
	}

	emailAlreadyTaken, err := h.userReader.UserWithEmailExists(r.Context(), req.NewEmail)
	if err != nil {
		h.logger.Error("failed to determine if the email address is taken", "error", err)
		send.InternalServerError(w, "There has been an issue updating your email address")
		return
	}
	if emailAlreadyTaken {
		send.BadRequest(w, "Email address already in use")
		return
	}

	result, err := h.initializeUpdateEmailCommand.Execute(r.Context(), c.Subject, req.NewEmail)
	if err != nil {
		h.logger.Error("failed to execute initialize update email command", "error", err)
		send.InternalServerError(w, "There has been an issue updating your email address")
		return
	}

	h.sendUpdateEmailAddressEmail(c, req.NewEmail, result.OTP)
	w.WriteHeader(http.StatusOK)
}

func (h *InitializeEmailUpdateHandler) sendUpdateEmailAddressEmail(c claims.JWTClaims, newEmail, otp string) {
	html, err := email.GenerateUpdateEmailAddressEmail(email.UpdateEmailAddressData{
		Username:          c.Username,
		NewEmail:          newEmail,
		OTP:               otp,
		ExpirationMinutes: strconv.FormatInt(int64(command.OTPTimeToLiveMinutes), 10),
	})
	if err != nil {
		h.logger.Error("failed to generate update email address email template", "error", err)
		return
	}
	if _, err := h.emailSender.Send(c.Email, "Update email address", html); err != nil {
		h.logger.Error("failed to send email", "error", err)
	}
}
