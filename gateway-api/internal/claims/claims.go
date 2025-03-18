package claims

import (
	"net/http"
	"time"
)

const (
	JWTClaimsKey    = "claims"
	RefreshTokenKey = "refresh_token"
	AccessTokenKey  = "access_token"
)

type JWTClaims struct {
	Expiration int64  `json:"exp"`
	Subject    string `json:"sub"`
	Username   string `json:"username"`
}

func (c JWTClaims) Authenticated() bool {
	return c.Subject != "" && c.Username != "" && !c.Expired()
}

func (c JWTClaims) Expired() bool {
	return time.Unix(c.Expiration, 0).Before(time.Now())
}

func GetClaims(r *http.Request) JWTClaims {
	claims, ok := r.Context().Value(JWTClaimsKey).(JWTClaims)
	if !ok {
		return JWTClaims{}
	}
	return claims
}

func GetRefreshToken(r *http.Request) (string, bool) {
	token, ok := r.Context().Value(RefreshTokenKey).(string)
	return token, ok
}
