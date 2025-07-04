package middleware

import (
	"beerbux/internal/auth/command"
	"beerbux/internal/auth/cookie"
	"beerbux/internal/common/claims"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
)

type AuthMiddleware struct {
	refreshTokenCommand *command.RefreshTokenCommand
	secret              string
}

func NewAuthMiddleware(refreshTokenCommand *command.RefreshTokenCommand, secret string) *AuthMiddleware {
	return &AuthMiddleware{
		refreshTokenCommand: refreshTokenCommand,
		secret:              secret,
	}
}

// WithJWT is a middleware that extracts the JWT claims from the request and adds them to the context.
// If the JWT is invalid or does not exist, the middleware will continue to the next handler.
func (mw *AuthMiddleware) WithJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		accessCookie, err := r.Cookie(cookie.AccessTokenKey)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		refreshCookie, err := r.Cookie(cookie.RefreshTokenKey)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		accessToken := accessCookie.Value
		refreshToken := refreshCookie.Value

		var ve *jwt.ValidationError
		jwtClaims, err := mw.parseJWTClaims(accessToken)
		if err != nil && errors.As(err, &ve) && ve.Errors == jwt.ValidationErrorExpired {
			// Attempt to refresh the access token if it has expired.
			user, err := mw.refreshAccessToken(r.Context(), accessToken, refreshToken)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			jwtClaims, err = mw.parseJWTClaims(user.AccessToken)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			accessToken = user.AccessToken
			refreshToken = user.RefreshToken
			cookie.SetAccessTokenCookie(w, accessToken)
			cookie.SetRefreshTokenCookie(w, refreshToken)
		} else if err != nil {
			// If the JWT is invalid for any other reason.
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), claims.JWTClaimsKey, jwtClaims)
		ctx = context.WithValue(ctx, claims.RefreshTokenKey, refreshToken)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (mw *AuthMiddleware) refreshAccessToken(ctx context.Context, accessToken, refreshToken string) (*command.TokenResponse, error) {
	if accessToken == "" || refreshToken == "" {
		return nil, errors.New("missing access or refresh token")
	}

	subject, err := mw.unsafeGetSubjectFromJWT(accessToken)
	if err != nil {
		return nil, err
	}

	return mw.refreshTokenCommand.Execute(ctx, subject, refreshToken)
}

func (mw *AuthMiddleware) parseJWTClaims(jwtValue string) (claims.JWTClaims, error) {
	token, err := jwt.Parse(jwtValue, func(token *jwt.Token) (interface{}, error) {
		return []byte(mw.secret), nil
	})
	if err != nil {
		return claims.JWTClaims{}, fmt.Errorf("failed parsing access token: %w", err)
	}

	if !token.Valid {
		return claims.JWTClaims{}, fmt.Errorf("invalid access token")
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return claims.JWTClaims{}, fmt.Errorf("invalid claims")
	}

	return mw.parseClaimValues(mapClaims)
}

// unsafeGetSubjectFromJWT returns the subject from the JWT without verifying the signature.
// This is intended for determining the subject when refreshing the JWT.
func (mw *AuthMiddleware) unsafeGetSubjectFromJWT(jwtValue string) (uuid.UUID, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(jwtValue, jwt.MapClaims{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed parsing access token: %w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims")
	}

	subValue, ok := mapClaims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("sub claim not found")
	}

	return uuid.Parse(subValue)
}

func (mw *AuthMiddleware) parseClaimValues(mapClaims jwt.MapClaims) (claims.JWTClaims, error) {
	subValue, ok := mapClaims["sub"].(string)
	if !ok {
		return claims.JWTClaims{}, fmt.Errorf("sub claim not found")
	}
	if _, err := uuid.Parse(subValue); err != nil {
		return claims.JWTClaims{}, fmt.Errorf("invalid user id")
	}

	username, ok := mapClaims["username"].(string)
	if !ok {
		return claims.JWTClaims{}, fmt.Errorf("sub claim not found")
	}

	email, ok := mapClaims["email"].(string)
	if !ok {
		return claims.JWTClaims{}, fmt.Errorf("email claim not found")
	}

	exp, ok := mapClaims["exp"].(float64)
	if !ok {
		return claims.JWTClaims{}, fmt.Errorf("sub claim not found")
	}

	sub, err := uuid.Parse(subValue)
	if err != nil {
		return claims.JWTClaims{}, err
	}

	return claims.JWTClaims{
		Subject:    sub,
		Username:   username,
		Email:      email,
		Expiration: int64(exp),
	}, nil
}
