package server_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/auth"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/token"
	"github.com/thisisthemurph/beerbux/auth-service/internal/server"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/auth-service/tests/builder"
	"github.com/thisisthemurph/beerbux/auth-service/tests/testinfra"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "supersecret"

func setupAuthServer(db *sql.DB) *server.AuthServer {
	authRepo := auth.New(db)
	authTokenRepo := token.New(db)
	return server.NewAuthServer(authRepo, authTokenRepo, server.AuthServerOptions{
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  time.Hour,
		RefreshTokenTTL: time.Hour * 24 * 7,
	})
}

func TestAuthServer_Login(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	authServer := setupAuthServer(db)

	user := builder.NewUserBuilder(t).
		WithUsername("testuser").
		WithPassword("password").
		Build(db)

	testCases := []struct {
		name          string
		username      string
		password      string
		expectedError error
	}{
		{
			"valid credentials",
			user.Username,
			"password",
			nil,
		},
		{
			"invalid password",
			user.Username,
			"wrong-password",
			server.ErrInvalidCredentials,
		},
		{
			"invalid username",
			"wrong.user",
			"password",
			server.ErrInvalidCredentials,
		},
		{
			"invalid username and password",
			"wrong.user",
			"wrong-password",
			server.ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := authServer.Login(context.Background(), &authpb.LoginRequest{
				Username: tc.username,
				Password: tc.password,
			})

			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.AccessToken)
				assert.Len(t, res.RefreshToken, 44)
				validateAuthToken(t, res.AccessToken, tc.username)
			}
		})
	}
}

func TestAuthServer_Signup_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	authServer := setupAuthServer(db)

	req := authpb.SignupRequest{
		Name:                 "name",
		Username:             "username",
		Password:             "password",
		VerificationPassword: "password",
	}

	res, err := authServer.Signup(context.Background(), &req)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	var username, hashedPassword string
	q := "select username, hashed_password from users where username = ? limit 1;"
	err = db.QueryRow(q, req.Username).Scan(&username, &hashedPassword)
	require.NoError(t, err)
	assert.Equal(t, req.Username, username)
	assert.NotEmpty(t, hashedPassword)
	validPassword(t, req.Password, hashedPassword)
}

func TestAuthServer_Signup_WithBadRequests(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	authServer := setupAuthServer(db)

	existingUser := builder.NewUserBuilder(t).
		WithUsername("existing").
		WithPassword("password").
		Build(db)

	testCases := []struct {
		name          string
		req           *authpb.SignupRequest
		expectedError error
	}{
		{
			"password mismatch",
			&authpb.SignupRequest{
				Name:                 "name",
				Username:             "username",
				Password:             "password",
				VerificationPassword: "password2",
			},
			server.ErrPasswordsDoNotMatch,
		},
		{
			"existing username",
			&authpb.SignupRequest{
				Name:                 "name",
				Username:             existingUser.Username,
				Password:             "password",
				VerificationPassword: "password",
			},
			server.ErrUsernameExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := authServer.Signup(context.Background(), tc.req)
			assert.Error(t, err)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Nil(t, res)
		})
	}
}

func TestAuthServer_RefreshToken(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	authServer := setupAuthServer(db)

	user := builder.NewUserBuilder(t).
		WithUsername("testuser").
		WithPassword("password").
		Build(db)

	tokenRaw := "hashed-token"
	hashedToken, _ := bcrypt.GenerateFromPassword([]byte(tokenRaw), bcrypt.DefaultCost)
	existingRefreshToken := builder.NewRefreshTokenBuilder(t).
		WithUserID(user.ID).
		WithHashedToken(string(hashedToken)).
		WithExpiresAt(time.Now().Add(time.Hour)).
		Build(db)

	testCases := []struct {
		name          string
		userID        string
		refreshToken  string
		expectedError error
	}{
		{
			"valid refresh token",
			user.ID,
			tokenRaw,
			nil,
		},
		{
			"unexpected refresh token",
			user.ID,
			"unexpected-token",
			server.ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := authServer.RefreshToken(context.Background(), &authpb.RefreshTokenRequest{
				UserId:       tc.userID,
				RefreshToken: tc.refreshToken,
			})

			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.AccessToken)
				assert.Len(t, res.RefreshToken, 44)
				assert.NotEqual(t, tc.refreshToken, res.RefreshToken)
				validateAuthToken(t, res.AccessToken, user.Username)

				// Check that the old refresh token is deleted.
				q := "select exists(select 1 from refresh_tokens where id = ?);"
				var exists bool
				err = db.QueryRow(q, existingRefreshToken.ID).Scan(&exists)
				require.NoError(t, err)
				assert.False(t, exists)
			} else {
				assert.Nil(t, res)
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedError)
			}
		})
	}
}

func validPassword(t *testing.T, expected, actual string) {
	err := bcrypt.CompareHashAndPassword([]byte(actual), []byte(expected))
	assert.NoError(t, err)
}

func validateAuthToken(t *testing.T, tokenString, expectedUsername string) {
	accessToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Fatal("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		t.Error(err)
	}

	assert.True(t, accessToken.Valid)
	assert.NoError(t, accessToken.Claims.Valid())

	claims, ok := accessToken.Claims.(jwt.MapClaims)
	assert.True(t, ok, "token.Claims is not jwt.MapClaims")
	assert.NotEmpty(t, claims["username"])
	assert.NotEmpty(t, claims["exp"])

	assert.Equal(t, expectedUsername, claims["username"])
}
