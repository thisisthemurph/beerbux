package handler

import (
	"beerbux/internal/common/claims"
	"beerbux/internal/session/command"
	"beerbux/pkg/send"
	"encoding/json"
	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type CreateSessionHandler struct {
	createSessionCommand *command.CreateSessionCommand
	logger               *slog.Logger
}

func NewCreateSessionHandler(createSessionCommand *command.CreateSessionCommand, logger *slog.Logger) *CreateSessionHandler {
	return &CreateSessionHandler{
		createSessionCommand: createSessionCommand,
		logger:               logger,
	}
}

type CreateSessionRequest struct {
	Name string `json:"name"`
}

type CreateSessionResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (h *CreateSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		send.ValidationError(w, err)
		return
	}

	s, err := h.createSessionCommand.Execute(r.Context(), c.Subject, req.Name)
	if err != nil {
		h.logger.Error("error creating session", "error", err)
		send.InternalServerError(w, "Failed to create session")
		return
	}

	send.JSON(w, CreateSessionResponse{
		ID:   s.ID,
		Name: s.Name,
	}, http.StatusCreated)
}

func (r CreateSessionRequest) Validate() error {
	return oz.ValidateStruct(&r,
		oz.Field(&r.Name,
			oz.Required.Error("A name must be provided"),
			oz.Length(3, 25).Error("The name must be between 2 and 25 characters"),
		),
	)
}
