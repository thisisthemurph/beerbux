package server_test

import (
	"context"
	"database/sql"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/ledger"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/beerbux/user-service/internal/repository/user"
	"github.com/thisisthemurph/beerbux/user-service/internal/server"
	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"github.com/thisisthemurph/beerbux/user-service/tests/builder"
	"github.com/thisisthemurph/beerbux/user-service/tests/fake"
	"github.com/thisisthemurph/beerbux/user-service/tests/testinfra"
)

func setupUserServer(db *sql.DB) *server.UserServer {
	userRepo := user.New(db)
	userLedgerRepo := ledger.New(db)
	fakeUserCreatedPublished := fake.NewFakeUserCreatedPublisher()
	fakeUserUpdatedPublished := fake.NewFakeUserUpdatedPublisher()
	return server.NewUserServer(userRepo, userLedgerRepo, fakeUserCreatedPublished, fakeUserUpdatedPublished, slog.Default())
}

func TestGetUser_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	userServer := setupUserServer(db)

	usr := builder.NewUserBuilder(t).
		WithName("Test User").
		WithUsername("testuser").
		Build(db)

	res, err := userServer.GetUser(context.Background(), &userpb.GetUserRequest{UserId: usr.ID})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, usr.ID, res.UserId)
	assert.Equal(t, usr.Name, res.Name)
	assert.Equal(t, usr.Username, res.Username)
}

func TestCreateUser_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	userServer := setupUserServer(db)

	bio := "This is a test user"
	res, err := userServer.CreateUser(context.Background(), &userpb.CreateUserRequest{
		Name:     "Test User",
		Username: "testuser",
		Bio:      &bio,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.UserId)
	assert.Equal(t, "Test User", res.Name)
	assert.Equal(t, "testuser", res.Username)
	assert.Equal(t, bio, *res.Bio)
}

func TestUpdateUser_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	userServer := setupUserServer(db)

	usr := builder.NewUserBuilder(t).
		WithName("Test User").
		WithUsername("testuser").
		Build(db)

	bio := "This is an updated test user"
	res, err := userServer.UpdateUser(context.Background(), &userpb.UpdateUserRequest{
		UserId:   usr.ID,
		Name:     "Updated User",
		Username: "updateduser",
		Bio:      &bio,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, usr.ID, res.UserId)
	assert.Equal(t, "Updated User", res.Name)
	assert.Equal(t, "updateduser", res.Username)

	var name, username string
	err = db.QueryRow("select name, username from users where id = ?", usr.ID).Scan(&name, &username)
	assert.NoError(t, err)
	assert.Equal(t, "Updated User", name)
	assert.Equal(t, "updateduser", username)
}
