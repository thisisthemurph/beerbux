package server_test

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log/slog"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/auth"
	"github.com/thisisthemurph/beerbux/auth-service/internal/repository/token"
	"github.com/thisisthemurph/beerbux/auth-service/internal/server"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"github.com/thisisthemurph/beerbux/auth-service/tests/builder"
	"github.com/thisisthemurph/beerbux/auth-service/tests/fakes"
	"github.com/thisisthemurph/beerbux/auth-service/tests/testinfra"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "supersecret"

func setupAuthServer(db *sql.DB) *server.AuthServer {
	authRepo := auth.New(db)
	authTokenRepo := token.New(db)
	userRegisteredProducer := fakes.NewFakeUserRegisteredProducer()
	return server.NewAuthServer(slog.Default(), authRepo, authTokenRepo, userRegisteredProducer, server.AuthServerOptions{
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
			server.ErrUserNotFound,
		},
		{
			"invalid username",
			"wrong.user",
			"password",
			server.ErrUserNotFound,
		},
		{
			"invalid username and password",
			"wrong.user",
			"wrong-password",
			server.ErrUserNotFound,
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
			server.ErrPasswordMismatch,
		},
		{
			"existing username",
			&authpb.SignupRequest{
				Name:                 "name",
				Username:             existingUser.Username,
				Password:             "password",
				VerificationPassword: "password",
			},
			server.ErrUsernameTaken,
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
			server.ErrRefreshTokenNotFound,
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

func TestAuthServer_InvalidateRefreshToken(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	authServer := setupAuthServer(db)

	user := builder.NewUserBuilder(t).
		WithUsername("testuser").
		WithPassword("password").
		Build(db)

	// A separate user to ensure other user's refresh token is not affected.
	user2 := builder.NewUserBuilder(t).
		WithUsername("user2").
		WithPassword("password").
		Build(db)

	tokenRaw := uuid.NewString()
	u1ht1, _ := bcrypt.GenerateFromPassword([]byte(tokenRaw), bcrypt.DefaultCost)
	user1Token1 := builder.NewRefreshTokenBuilder(t).
		WithUserID(user.ID).
		WithHashedToken(string(u1ht1)).
		WithExpiresAt(time.Now().Add(time.Hour)).
		Build(db)

	// Another refresh token for the same user, should not be affected.
	u1ht2, _ := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), bcrypt.DefaultCost)
	user1Token2 := builder.NewRefreshTokenBuilder(t).
		WithUserID(user.ID).
		WithHashedToken(string(u1ht2)).
		WithExpiresAt(time.Now().Add(time.Hour)).
		Build(db)

	// A refresh token for user2, should not be affected.
	u2ht, _ := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), bcrypt.DefaultCost)
	user2Token := builder.NewRefreshTokenBuilder(t).
		WithUserID(user2.ID).
		WithHashedToken(string(u2ht)).
		WithExpiresAt(time.Now().Add(time.Hour)).
		Build(db)

	testCases := []struct {
		name          string
		userID        string
		refreshToken  string
		expectedError error
	}{
		{
			"existing refresh token",
			user.ID,
			tokenRaw,
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := authServer.InvalidateRefreshToken(context.Background(), &authpb.InvalidateRefreshTokenRequest{
				UserId:       tc.userID,
				RefreshToken: tc.refreshToken,
			})

			assert.NoError(t, err)

			var revoked bool
			q := "select revoked from refresh_tokens where id = ?;"

			// The main user's first token should be revoked.
			err = db.QueryRow(q, user1Token1.ID).Scan(&revoked)
			require.NoError(t, err)
			assert.True(t, revoked)

			// The main user's second token should not be revoked.
			err = db.QueryRow(q, user1Token2.ID).Scan(&revoked)
			require.NoError(t, err)
			assert.False(t, revoked)

			// The other user's token should not be revoked.
			err = db.QueryRow(q, user2Token.ID).Scan(&revoked)
			require.NoError(t, err)
			assert.False(t, revoked)

			// There should be 3 refresh tokens in the database.
			q = "select count(*) from refresh_tokens;"
			var count int
			err = db.QueryRow(q).Scan(&count)
			require.NoError(t, err)
			assert.Equal(t, 3, count)
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
