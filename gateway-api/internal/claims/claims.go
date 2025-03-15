package claims

import "net/http"

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

func GetClaims(r *http.Request) (*JWTClaims, bool) {
	claims, ok := r.Context().Value(JWTClaimsKey).(*JWTClaims)
	return claims, ok
}

func GetSubject(r *http.Request) (string, bool) {
	claims, ok := GetClaims(r)
	if !ok {
		return "", false
	}
	return claims.Subject, true
}

func GetUsername(r *http.Request) (string, bool) {
	claims, ok := GetClaims(r)
	if !ok {
		return "", false
	}
	return claims.Username, true
}

func GetRefreshToken(r *http.Request) (string, bool) {
	token, ok := r.Context().Value(RefreshTokenKey).(string)
	return token, ok
}
