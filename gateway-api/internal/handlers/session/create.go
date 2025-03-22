package session

import (
	"encoding/json"
	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"net/http"
)

type CreateSessionHandler struct {
	sessionClient sessionpb.SessionClient
}

func NewCreateSessionHandler(sessionClient sessionpb.SessionClient) *CreateSessionHandler {
	return &CreateSessionHandler{
		sessionClient: sessionClient,
	}
}

type CreateSessionRequest struct {
	Name string `json:"name"`
}

type CreateSessionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

	ssn, err := h.sessionClient.CreateSession(r.Context(), &sessionpb.CreateSessionRequest{
		UserId: c.Subject,
		Name:   req.Name,
	})

	if err != nil {
		send.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	send.JSON(w, CreateSessionResponse{
		ID:   ssn.SessionId,
		Name: ssn.Name,
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
