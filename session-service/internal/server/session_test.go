package server_test

import (
	"context"
	"database/sql"
	"fmt"
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
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, slog.Default())

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

	var sessionName string
	err = db.QueryRow("select name from sessions where id = ?", resp.SessionId).Scan(&sessionName)
	assert.NoError(t, err)
	assert.Equal(t, "New Session", sessionName)

	assertUserInsertedAsMember(t, db, resp.SessionId, fakeUserID, "user", "username", true)
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

func TestAddMemberToSession_WhenMemberInMembersTable_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })
	sessionRepo := session.New(db)

	existingMember := builder.NewMemberBuilder(t).
		WithName("user").
		WithUsername("username").
		Build(db)

	ssn := builder.NewSessionBuilder(t).
		WithName("Test Session").
		Build(db)

	fakeUserClient := fake.NewFakeUserClient()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, slog.Default())

	req := &sessionpb.AddMemberToSessionRequest{
		SessionId: ssn.ID,
		UserId:    existingMember.ID,
	}

	_, err := sessionServer.AddMemberToSession(context.Background(), req)
	assert.NoError(t, err)
	assertUserInsertedAsMember(t, db, ssn.ID, existingMember.ID, "user", "username", false)
}

func TestAddMemberToSession_WhenMemberNotInSessionsTable_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })
	sessionRepo := session.New(db)

	testUserID := uuid.NewString()
	ssn := builder.NewSessionBuilder(t).
		WithName("Test Session").
		Build(db)

	fakeUserClient := fake.NewFakeUserClient().WithUser(testUserID, "user", "username")
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, slog.Default())

	req := &sessionpb.AddMemberToSessionRequest{
		SessionId: ssn.ID,
		UserId:    testUserID,
	}

	_, err := sessionServer.AddMemberToSession(context.Background(), req)
	assert.NoError(t, err)
	assertUserInsertedAsMember(t, db, ssn.ID, testUserID, "user", "username", false)
}

func TestAddMemberToSession_WhenSessionNotFound_Errors(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })
	sessionRepo := session.New(db)

	testUserID := uuid.NewString()
	fakeUserClient := fake.NewFakeUserClient().WithUser(testUserID, "user", "username")
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, slog.Default())

	req := &sessionpb.AddMemberToSessionRequest{
		SessionId: uuid.NewString(),
		UserId:    testUserID,
	}

	_, err := sessionServer.AddMemberToSession(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("failed fetching session %q from database", req.SessionId))
	assert.Contains(t, err.Error(), "sql: no rows in result set")
}

func TestAddMemberToSession_WhenUserNotFound_Errors(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })
	sessionRepo := session.New(db)

	ssn := builder.NewSessionBuilder(t).
		WithName("Test Session").
		Build(db)

	testUserID := uuid.NewString()
	fakeUserClient := fake.NewFakeUserClient() //.WithUser(testUserID, "user", "username")
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, slog.Default())

	req := &sessionpb.AddMemberToSessionRequest{
		SessionId: ssn.ID,
		UserId:    testUserID,
	}

	_, err := sessionServer.AddMemberToSession(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch user: user not found")
}

func assertUserInsertedAsMember(t *testing.T, db *sql.DB, sessionID, memberID, expectedName, expectedUsername string, expectedOwner bool) {
	var memberName, memberUsername string
	err := db.QueryRow("select name, username from members where id = ?", memberID).Scan(&memberName, &memberUsername)
	assert.NoError(t, err)
	assert.Equal(t, expectedName, memberName)
	assert.Equal(t, expectedUsername, memberUsername)

	var isOwner bool
	err = db.QueryRow("select is_owner from session_members where session_id = ? and member_id = ?", sessionID, memberID).Scan(&isOwner)
	assert.NoError(t, err)
	assert.Equal(t, expectedOwner, isOwner)
}
