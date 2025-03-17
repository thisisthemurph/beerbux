package auth

import (
	"encoding/json"
	"net/http"

	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers"
)

type LoginHandler struct {
	authClient authpb.AuthClient
}

func NewLoginHandler(authClient authpb.AuthClient) *LoginHandler {
	return &LoginHandler{
		authClient: authClient,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		handlers.WriteValidationError(w, err)
		return
	}

	loginRequest := authpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	loginResp, err := h.authClient.Login(r.Context(), &loginRequest)

	if err != nil {
		handlers.WriteValidationError(w, err)
		return
	}

	setAccessTokenCookie(w, loginResp.AccessToken)
	setRefreshTokenCookie(w, loginResp.RefreshToken)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		ID:       loginResp.User.Id,
		Username: loginResp.User.Username,
	})
}

func (r LoginRequest) Validate() error {
	return oz.ValidateStruct(&r,
		oz.Field(&r.Username, oz.Required),
		oz.Field(&r.Password, oz.Required),
	)
}
