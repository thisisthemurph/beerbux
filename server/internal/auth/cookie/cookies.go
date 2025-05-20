package cookie

import (
	"net/http"
	"time"
)

const (
	AccessTokenKey  = "access_token"
	RefreshTokenKey = "refresh_token"
)

func SetAccessTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenKey,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenKey,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(7 * time.Hour * 24),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}
