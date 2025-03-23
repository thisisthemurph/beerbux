package session

import (
	"encoding/json"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"net/http"
)

type AddMemberToSessionHandler struct {
	userClient    userpb.UserClient
	sessionClient sessionpb.SessionClient
}

func NewAddMemberToSessionHandler(userClient userpb.UserClient, sessionClient sessionpb.SessionClient) *AddMemberToSessionHandler {
	return &AddMemberToSessionHandler{
		userClient:    userClient,
		sessionClient: sessionClient,
	}
}

type AddMemberToSessionRequest struct {
	Username string `json:"username"`
}

func (h *AddMemberToSessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := claims.GetClaims(r)
	if !c.Authenticated() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionID, err := uuid.Parse(r.PathValue("sessionId"))
	if err != nil {
		send.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	var req AddMemberToSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		send.ValidationError(w, err)
		return
	}

	ssn, err := h.sessionClient.GetSession(r.Context(), &sessionpb.GetSessionRequest{
		SessionId: sessionID.String(),
	})
	if err != nil {
		send.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	// Ensure the current user is a member of the session.
	currentUserIsMember := false
	for _, m := range ssn.Members {
		if m.UserId == c.Subject {
			currentUserIsMember = true
			break
		}
	}
	if !currentUserIsMember {
		send.Error(w, "You do not have permission to add members to this session", http.StatusUnauthorized)
		return
	}

	//Fetch the user being added to the session
	userToAdd, err := h.userClient.GetUserByUsername(r.Context(), &userpb.GetUserByUsernameRequest{
		Username: req.Username,
	})
	if err != nil {
		sc, ok := status.FromError(err)
		if ok && sc.Code() == codes.NotFound {
			send.Error(w, "User not found", http.StatusNotFound)
			return
		}
		send.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if the new user is already a member of the session.
	for _, m := range ssn.Members {
		if m.UserId == userToAdd.UserId {
			// Exit early with a 200 rather than a 201 to indicate that the member was already added.
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	_, err = h.sessionClient.AddMemberToSession(r.Context(), &sessionpb.AddMemberToSessionRequest{
		SessionId: sessionID.String(),
		UserId:    userToAdd.UserId,
	})
	if err != nil {
		send.Error(w, "Failed to add member to session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (r AddMemberToSessionRequest) Validate() error {
	return oz.ValidateStruct(&r,
		oz.Field(&r.Username,
			oz.Required.Error("The member's username is required"),
		),
	)
}
