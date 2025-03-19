package auth

import (
	"encoding/json"
	"net/http"

	oz "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/cookie"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/handlers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		st, ok := status.FromError(err)
		if !ok {
			handlers.WriteError(w, "There has been an error logging you in", http.StatusUnauthorized)
			return
		}

		switch st.Code() {
		case codes.NotFound:
			handlers.WriteError(w, "Invalid username or password", http.StatusUnauthorized)
			return
		default:
			handlers.WriteError(w, "There has been an error logging you in", http.StatusUnauthorized)
			return
		}
	}

	cookie.SetAccessTokenCookie(w, loginResp.AccessToken)
	cookie.SetRefreshTokenCookie(w, loginResp.RefreshToken)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(LoginResponse{
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
