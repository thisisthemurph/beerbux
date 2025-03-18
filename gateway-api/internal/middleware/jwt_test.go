package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/claims"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/cookie"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/middleware"
	"github.com/thisisthemurph/beerbux/gateway-api/internal/tests/fake"
)

const secret = "thisisasecretstringforsigningjwttokens"

func TestWithJWT_WithValidJWTInCookie_SetsValuesInContext(t *testing.T) {
	sub := uuid.NewString()
	username := "username"
	expiration := time.Now().Add(time.Hour * 1).Unix()
	refreshToken := uuid.NewString()
	accessToken := generateAccessToken(t, secret, sub, username, expiration)

	fakeAuthClient := fake.NewFakeAuthClient(nil)
	authMw := middleware.NewAuthMiddleware(fakeAuthClient, secret)

	nextHandler, capturedCtx := captureContext(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	handler := authMw.WithJWT(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(newCookie(cookie.AccessTokenKey, accessToken))
	req.AddCookie(newCookie(cookie.RefreshTokenKey, refreshToken))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Ensure no additional cookies were set
	assert.Empty(t, rr.Result().Cookies())

	// Validate context claims
	require.NotNil(t, *capturedCtx, "Context should be populated by middleware")
	expectedClaims := claims.JWTClaims{
		Subject:    sub,
		Username:   username,
		Expiration: expiration,
	}
	actualClaims, ok := (*capturedCtx).Value(claims.JWTClaimsKey).(claims.JWTClaims)
	assert.True(t, ok, "JWTClaims should be present in context")
	assert.Equal(t, expectedClaims, actualClaims)

	assert.Equal(t, refreshToken, (*capturedCtx).Value(claims.RefreshTokenKey))
}

func TestWithJWT_WithExpiredJWTInCookie_SetsCookiesAndValuesInContext(t *testing.T) {
	sub := uuid.NewString()
	username := "username"
	refreshToken := uuid.NewString()
	expiration := time.Now().Add(time.Hour * 1).Unix()
	newAccessToken := generateAccessToken(t, secret, sub, username, expiration)

	fakeAuthClient := fake.NewFakeAuthClient(&authpb.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	})

	authMw := middleware.NewAuthMiddleware(fakeAuthClient, secret)

	nextHandler, capturedCtx := captureContext(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	handler := authMw.WithJWT(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(newCookie(cookie.AccessTokenKey, generateExpiredAccessToken(t, secret, sub, username)))
	req.AddCookie(newCookie(cookie.RefreshTokenKey, refreshToken))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Ensure refreshed tokens are set
	cookies := rr.Result().Cookies()
	require.Len(t, cookies, 2)

	// Validate cookie values
	cookieMap := make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	assert.Equal(t, newAccessToken, cookieMap[cookie.AccessTokenKey])
	assert.Equal(t, refreshToken, cookieMap[cookie.RefreshTokenKey])

	// Validate context claims
	require.NotNil(t, *capturedCtx, "Context should be populated by middleware")
	expectedClaims := claims.JWTClaims{
		Subject:    sub,
		Username:   username,
		Expiration: expiration,
	}
	actualClaims, ok := (*capturedCtx).Value(claims.JWTClaimsKey).(claims.JWTClaims)
	require.True(t, ok, "JWTClaims should be present in context")
	assert.Equal(t, expectedClaims, actualClaims)

	assert.Equal(t, refreshToken, (*capturedCtx).Value(claims.RefreshTokenKey))
}

func TestWithJWT_WithNoTokenInCookie_DoesNotSetValuesInContext(t *testing.T) {
	fakeAuthClient := fake.NewFakeAuthClient(nil)
	authMw := middleware.NewAuthMiddleware(fakeAuthClient, secret)

	nextHandler, capturedCtx := captureContext(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	handler := authMw.WithJWT(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Ensure no additional cookies were set
	assert.Empty(t, rr.Result().Cookies())

	// Validate context claims
	require.NotNil(t, *capturedCtx, "Context should be populated by middleware")
	assert.Nil(t, (*capturedCtx).Value(claims.JWTClaimsKey))
	assert.Nil(t, (*capturedCtx).Value(claims.RefreshTokenKey))
}

func TestWithJWT_WithEmptyRefreshToken_DoesNotSetValuesOnContext(t *testing.T) {
	sub := uuid.NewString()
	username := "username"
	expiration := time.Now().Add(time.Hour * 1).Unix()
	accessToken := generateAccessToken(t, secret, sub, username, expiration)

	fakeAuthClient := fake.NewFakeAuthClient(nil)
	authMw := middleware.NewAuthMiddleware(fakeAuthClient, secret)

	nextHandler, capturedCtx := captureContext(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	handler := authMw.WithJWT(nextHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(newCookie(cookie.AccessTokenKey, accessToken))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Ensure no cookies were set
	assert.Empty(t, rr.Result().Cookies())

	// Validate context claims
	require.NotNil(t, *capturedCtx, "Context should be populated by middleware")
	assert.Nil(t, (*capturedCtx).Value(claims.JWTClaimsKey))
	assert.Nil(t, (*capturedCtx).Value(claims.RefreshTokenKey))
}

func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

func captureContext(next http.HandlerFunc) (http.HandlerFunc, *context.Context) {
	var capturedCtx context.Context
	return func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		next(w, r)
	}, &capturedCtx
}

func generateAccessToken(t *testing.T, secret, sub, username string, expiration int64) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      sub,
		"username": username,
		"exp":      expiration,
	})

	ts, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return ts
}

func generateExpiredAccessToken(t *testing.T, secret, sub, username string) string {
	return generateAccessToken(t, secret, sub, username, time.Now().Add(-1*time.Hour).Unix())
}
