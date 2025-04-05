package auth

import (
	"encoding/json"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/dto"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"net/http"

	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/cookie"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers/shared/send"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginHandler struct {
	authClient authpb.AuthClient
	userClient userpb.UserClient
}

func NewLoginHandler(authClient authpb.AuthClient, userClient userpb.UserClient) *LoginHandler {
	return &LoginHandler{
		authClient: authClient,
		userClient: userClient,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		send.Error(w, "Failed to decode request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		send.ValidationError(w, err)
		return
	}

	// Authenticate the user via the auth service
	loginResp, err := h.authClient.Login(r.Context(), &authpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		h.handleLoginError(w, err)
		return
	}

	// Get the user details from the user service
	user, err := h.userClient.GetUser(r.Context(), &userpb.GetUserRequest{
		UserId: loginResp.User.Id,
	})
	if err != nil {
		send.Error(w, "Failed to get user details", http.StatusInternalServerError)
		return
	}

	cookie.SetAccessTokenCookie(w, loginResp.AccessToken)
	cookie.SetRefreshTokenCookie(w, loginResp.RefreshToken)

	send.JSON(w, dto.UserResponse{
		ID:       user.UserId,
		Name:     user.Name,
		Username: user.Username,
	}, http.StatusOK)
}

func (h *LoginHandler) handleLoginError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		send.Error(w, "There has been an error logging you in", http.StatusUnauthorized)
		return
	}

	switch st.Code() {
	case codes.NotFound:
		send.Error(w, "Invalid username or password", http.StatusUnauthorized)
	default:
		send.Error(w, "There has been an error logging you in", http.StatusInternalServerError)
	}
}

func (r LoginRequest) Validate() error {
	return oz.ValidateStruct(&r,
		oz.Field(&r.Username, oz.Required),
		oz.Field(&r.Password, oz.Required),
	)
}
