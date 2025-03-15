package auth

import (
	"net/http"
	"time"
)

const (
	cookieAccessTokenKey  = "access_token"
	cookieRefreshTokenKey = "refresh_token"
)

func setAccessTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieAccessTokenKey,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieRefreshTokenKey,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(7 * time.Hour * 24),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}
