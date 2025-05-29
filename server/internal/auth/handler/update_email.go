package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/common/claims"
	"log/slog"
	"net/http"
)

type UpdateEmailHandler struct {
	updateEmailCommand *command.UpdateEmailCommand
	logger             *slog.Logger
}

func NewUpdateEmailHandler(updateEmailCommand *command.UpdateEmailCommand, logger *slog.Logger) *UpdateEmailHandler {
	return &UpdateEmailHandler{
		updateEmailCommand: updateEmailCommand,
		logger:             logger,
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
}
