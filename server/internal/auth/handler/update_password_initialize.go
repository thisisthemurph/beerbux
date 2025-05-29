package handler

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/common/claims"
	"beerbux/pkg/send"
	"encoding/json"
	"log/slog"
	"net/http"
)

type InitializePasswordUpdateHandler struct {
	initializePasswordUpdateCommand *command.InitializePasswordResetCommand
	logger                          *slog.Logger
	isDevelopmentMode               bool
}

func NewInitializePasswordUpdateHandler(
	initializePasswordResetCommand *command.InitializePasswordResetCommand,
	isDevelopmentMode bool,
	logger *slog.Logger,
) *InitializePasswordUpdateHandler {
	return &InitializePasswordUpdateHandler{
		initializePasswordUpdateCommand: initializePasswordResetCommand,
		isDevelopmentMode:               isDevelopmentMode,
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

	if h.isDevelopmentMode {
		h.logger.Info("Password reset initialized", "user", c.Subject, "otp", result.OTP)
	}

	w.WriteHeader(http.StatusOK)
}
