package server_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/internal/server"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/session-service/tests/builder"
	"github.com/thisisthemurph/beerbux/session-service/tests/fake"
	"github.com/thisisthemurph/beerbux/session-service/tests/testinfra"

	_ "modernc.org/sqlite"
)

func TestGetSession_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	defer db.Close()
	queries := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	sessionServer := server.NewSessionServer(db, queries, fakeUserClient, slog.Default())

	ssn := builder.NewSessionBuilder(t).
		WithName("Test Session").
		Build(db)

	resp, err := sessionServer.GetSession(context.Background(), &sessionpb.GetSessionRequest{SessionId: ssn.ID})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, ssn.ID, resp.SessionId)
	assert.Equal(t, ssn.Name, resp.Name)
	assert.True(t, resp.IsActive)
}

func TestGetSession_NotFound(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, slog.Default())

	sessionID := uuid.New().String()
	resp, err := sessionServer.GetSession(context.Background(), &sessionpb.GetSessionRequest{SessionId: sessionID})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to get session")
}

func TestGetSession_InvalidRequest(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, slog.Default())

	resp, err := sessionServer.GetSession(context.Background(), &sessionpb.GetSessionRequest{SessionId: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid request")
}

func TestCreateSession_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })
	sessionRepo := session.New(db)

	fakeUserID := uuid.NewString()
	fakeUserClient := fake.NewFakeUserClient().WithUser(fakeUserID, "user", "username")
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, slog.Default())

	req := &sessionpb.CreateSessionRequest{
		UserId: fakeUserID,
		Name:   "New Session",
	}
	resp, err := sessionServer.CreateSession(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "New Session", resp.Name)

	var count int
	err = db.QueryRow("select count(*) from sessions where id = ?", resp.SessionId).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	err = db.QueryRow("select count(*) from session_members where session_id = ?", resp.SessionId).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	var sessionName string
	err = db.QueryRow("select name from sessions where id = ?", resp.SessionId).Scan(&sessionName)
	assert.NoError(t, err)
	assert.Equal(t, "New Session", sessionName)

	var memberID string
	var memberIsOwner bool
	err = db.QueryRow("select member_id, is_owner from session_members where session_id = ?", resp.SessionId).Scan(&memberID, &memberIsOwner)
	assert.NoError(t, err)
	assert.Equal(t, fakeUserID, memberID)
	assert.True(t, memberIsOwner)

	var memberName, memberUsername string
	err = db.QueryRow("select name, username from members where id = ?", memberID).Scan(&memberName, &memberUsername)
	assert.NoError(t, err)
	assert.Equal(t, "user", memberName)
	assert.Equal(t, "username", memberUsername)
}

func TestCreateSession_WhenUserNotFound_Error(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })
	sessionRepo := session.New(db)

	fakeUserClient := fake.NewFakeUserClient().WithUser(uuid.NewString(), "user", "username")
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, slog.Default())

	req := &sessionpb.CreateSessionRequest{
		UserId: uuid.NewString(),
		Name:   "New Session",
	}

	resp, err := sessionServer.CreateSession(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to fetch user")
}
