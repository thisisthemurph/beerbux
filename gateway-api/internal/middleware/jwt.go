package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
)

// WithJWT is a middleware that extracts the JWT claims from the request and adds them to the context.
// If the JWT is invalid or does not exist, the middleware will continue to the next handler.
func WithJWT(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtClaims, err := parseClaimsFromJWT(r, secret)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), claims.JWTClaimsKey, jwtClaims)
		r = r.WithContext(ctx)

		refreshToken, err := parseRefreshToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx = context.WithValue(ctx, claims.RefreshTokenKey, refreshToken)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func parseClaimsFromJWT(r *http.Request, secret string) (*claims.JWTClaims, error) {
	accessCookie, err := r.Cookie(claims.AccessTokenKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(accessCookie.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return parseClaimValues(mapClaims)
}

func parseClaimValues(mapClaims jwt.MapClaims) (*claims.JWTClaims, error) {
	sub, ok := mapClaims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("sub claim not found")
	}
	if _, err := uuid.Parse(sub); err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	username, ok := mapClaims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("sub claim not found")
	}

	exp, ok := mapClaims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("sub claim not found")
	}

	return &claims.JWTClaims{
		Expiration: int64(exp),
		Subject:    sub,
		Username:   username,
	}, nil
}

func parseRefreshToken(r *http.Request) (string, error) {
	refreshCookie, err := r.Cookie(claims.RefreshTokenKey)
	if err != nil {
		return "", err
	}

	return refreshCookie.Value, nil
}
