package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/cookie"
)

type AuthMiddleware struct {
	authClient authpb.AuthClient
	secret     string
}

func NewAuthMiddleware(authClient authpb.AuthClient, secret string) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
		secret:     secret,
	}
}

// WithJWT is a middleware that extracts the JWT claims from the request and adds them to the context.
// If the JWT is invalid or does not exist, the middleware will continue to the next handler.
func (mw *AuthMiddleware) WithJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func (mw *AuthMiddleware) refreshAccessToken(ctx context.Context, accessToken, refreshToken string) (*authpb.RefreshTokenResponse, error) {
	if accessToken == "" || refreshToken == "" {
		return nil, errors.New("missing access or refresh token")
	}

	subject, err := mw.unsafeGetSubjectFromJWT(accessToken)
	if err != nil {
		return nil, err
	}

	user, err := mw.authClient.RefreshToken(ctx, &authpb.RefreshTokenRequest{
		UserId:       subject,
		RefreshToken: refreshToken,
	})

	return user, err
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
func (mw *AuthMiddleware) unsafeGetSubjectFromJWT(jwtValue string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(jwtValue, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed parsing access token: %w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	sub, ok := mapClaims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("sub claim not found")
	}

	return sub, nil
}

func (mw *AuthMiddleware) parseClaimValues(mapClaims jwt.MapClaims) (claims.JWTClaims, error) {
	sub, ok := mapClaims["sub"].(string)
	if !ok {
		return claims.JWTClaims{}, fmt.Errorf("sub claim not found")
	}
	if _, err := uuid.Parse(sub); err != nil {
		return claims.JWTClaims{}, fmt.Errorf("invalid user id")
	}

	username, ok := mapClaims["username"].(string)
	if !ok {
		return claims.JWTClaims{}, fmt.Errorf("sub claim not found")
	}

	exp, ok := mapClaims["exp"].(float64)
	if !ok {
		return claims.JWTClaims{}, fmt.Errorf("sub claim not found")
	}

	return claims.JWTClaims{
		Expiration: int64(exp),
		Subject:    sub,
		Username:   username,
	}, nil
}
