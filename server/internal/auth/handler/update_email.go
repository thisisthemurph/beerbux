package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/auth/cookie"
	"beerbux/internal/common/claims"
	"beerbux/pkg/send"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type UpdateEmailHandler struct {
	updateEmailCommand    *command.UpdateEmailCommand
	generateTokensCommand *command.GenerateTokensCommand
	logger                *slog.Logger
}

func NewUpdateEmailHandler(
	updateEmailCommand *command.UpdateEmailCommand,
	generateTokensCommand *command.GenerateTokensCommand,
	logger *slog.Logger,
) *UpdateEmailHandler {
	return &UpdateEmailHandler{
		updateEmailCommand:    updateEmailCommand,
		generateTokensCommand: generateTokensCommand,
		logger:                logger,
	}
}

type UpdateEmailRequest struct {
	OTP string `json:"otp"`
}

func (h *UpdateEmailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req UpdateEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.BadRequest(w, "Failed to decode request")
		return
	}

	if err := h.updateEmailCommand.Execute(r.Context(), c.Subject, req.OTP); err != nil {
		h.handleUpdateEmailError(w, err)
		return
	}

	tokens, err := h.generateTokensCommand.Execute(r.Context(), c.Username)
	if err != nil {
		send.InternalServerError(w, "There has been an issue re-authenticating you following updating your email address. Please try logging out and back in.")
		return
	}
	cookie.SetAccessTokenCookie(w, tokens.AccessToken)
	cookie.SetRefreshTokenCookie(w, tokens.RefreshToken)

	w.WriteHeader(http.StatusOK)
}

func (h *UpdateEmailHandler) handleUpdateEmailError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, command.ErrProcessNotInitialized):
		send.BadRequest(w, "An update for your email address was not requested")
	case errors.Is(err, command.ErrUserNotFound):
		send.NotFound(w, "User not found")
	case errors.Is(err, command.ErrIncorrectOTP):
		send.BadRequest(w, "The provided OTP is incorrect")
	case errors.Is(err, command.ErrOTPExpired):
		send.BadRequest(w, "Your OTP has expired, please start the process again")
	default:
		h.logger.Error("failed to update email address", "error", err)
		send.InternalServerError(w, "There has been an issue updating your email address")
	}
}
