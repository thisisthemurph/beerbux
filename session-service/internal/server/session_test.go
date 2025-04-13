package server_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thisisthemurph/beerbux/session-service/internal/repository/session"
	"github.com/thisisthemurph/beerbux/session-service/internal/server"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/session-service/tests/builder"
	"github.com/thisisthemurph/beerbux/session-service/tests/fake"
	"github.com/thisisthemurph/beerbux/session-service/tests/testinfra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetSession_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

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

func TestGetSession_GetSession_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

	sessionID := uuid.New()
	member1ID := uuid.NewString()
	member2ID := uuid.NewString()
	member3ID := uuid.NewString()
	transactionID := uuid.NewString()

	ssn := builder.NewSessionBuilder(t).
		WithID(sessionID).
		WithName("Test Session").
		WithMember(builder.SessionMemberParams{
			ID:       member1ID,
			Name:     "member1",
			Username: "username1",
			IsOwner:  true,
		}).
		WithMember(builder.SessionMemberParams{
			ID:       member2ID,
			Name:     "member2",
			Username: "username2",
		}).
		WithMember(builder.SessionMemberParams{
			ID:       member3ID,
			Name:     "member3",
			Username: "username3",
		}).
		WithTransaction(builder.SessionTransactionParams{
			ID:        transactionID,
			SessionID: sessionID.String(),
			CreatorID: member1ID,
			Lines: []builder.SessionTransactionLine{
				{MemberID: member2ID, Amount: 1},
				{MemberID: member3ID, Amount: 1},
			},
		}).
		Build(db)

	resp, err := sessionServer.GetSession(context.Background(), &sessionpb.GetSessionRequest{SessionId: ssn.ID})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, ssn.ID, resp.SessionId)
	assert.Equal(t, ssn.Name, resp.Name)
	assert.True(t, resp.IsActive)

	// Check the transactions
	transaction := resp.Transactions[0]

	assert.Len(t, resp.Transactions, 1)
	assert.Equal(t, transactionID, transaction.TransactionId)
	assert.Equal(t, member1ID, transaction.UserId)
	assert.NotEmpty(t, transaction.CreatedAt)
	assert.Len(t, transaction.Lines, 2)
	// Hack because we don't know the order of the transaction lines
	assert.True(t, transaction.Lines[0].UserId == member2ID || transaction.Lines[0].UserId == member3ID)
	assert.True(t, transaction.Lines[1].UserId == member2ID || transaction.Lines[1].UserId == member3ID)
	assert.NotEqual(t, transaction.Lines[0].UserId, transaction.Lines[1].UserId)

	// Check the members

	expectedMembers := map[string]string{
		"member1": "username1",
		"member2": "username2",
		"member3": "username3",
	}

	for _, m := range resp.Members {
		expectedUsername, exists := expectedMembers[m.Name]
		assert.True(t, exists)
		assert.Equal(t, expectedUsername, m.Username)
		delete(expectedMembers, m.Name)
	}

	// Ensure all expected members have been found
	assert.Len(t, expectedMembers, 0, "Not all members were checked")
}

func TestGetSession_When_SessionDoesNotExist_ReturnsError(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

	sessionID := uuid.New().String()
	resp, err := sessionServer.GetSession(context.Background(), &sessionpb.GetSessionRequest{SessionId: sessionID})
	assert.Error(t, err)
	assert.ErrorIs(t, err, server.ErrSessionNotFound)
	assert.Nil(t, resp)
}

func TestGetSession_With_InvalidRequest_ReturnsError(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

	resp, err := sessionServer.GetSession(context.Background(), &sessionpb.GetSessionRequest{SessionId: ""})
	assertStatusHasCode(t, err, codes.InvalidArgument)
	assert.Nil(t, resp)
}

func TestListSessionsForUser_ReturnsSessionsOrderedByUpdatedAt(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

	member := builder.NewMemberBuilder(t).
		WithName("Member 1").
		WithUsername("member1").
		Build(db)

	type sessionData struct {
		Name      string
		UpdatedAt time.Time
	}

	var sessions []sessionData
	baseTime := time.Now()

	const SessionsToMake = 50
	for i := 0; i < SessionsToMake; i++ {
		sessions = append(sessions, sessionData{
			Name:      fmt.Sprintf("session-%02d", i),
			UpdatedAt: baseTime.Add(-time.Duration(i) * time.Hour),
		})
	}

	rand.Shuffle(len(sessions), func(i, j int) {
		sessions[i], sessions[j] = sessions[j], sessions[i]
	})

	for _, s := range sessions {
		builder.NewSessionBuilder(t).
			WithName(s.Name).
			WithUpdatedAt(s.UpdatedAt).
			WithExistingMember(member).
			Build(db)
	}

	res, err := sessionServer.ListSessionsForUser(context.Background(), &sessionpb.ListSessionsForUserRequest{
		UserId: member.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res.Sessions, SessionsToMake)

	for i := 0; i < len(res.Sessions); i++ {
		current := res.Sessions[i]
		expectedName := fmt.Sprintf("session-%02d", i)
		assert.Equal(t, expectedName, current.Name)
	}
}

func TestListSessionsForUser_When_ProvidedWithPagingData_ReturnsPagedResults(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })

	sessionRepo := session.New(db)
	fakeUserClient := fake.NewFakeUserClient()
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

	member := builder.NewMemberBuilder(t).
		WithName("Member 1").
		WithUsername("member1").
		Build(db)

	sessionIdentifiers := make([]string, 50) // This slice should contain the ids of the items in the order returned.
	const SessionsToMake = 50
	for i := 0; i < SessionsToMake; i++ {
		ssn := builder.NewSessionBuilder(t).
			WithName(fmt.Sprintf("session-%02d", i)).
			// Adding the following updated time should ensure they are returned in
			// the order they are inserted; most recently updated first.
			WithUpdatedAt(time.Now().Add(time.Duration(SessionsToMake-i) * time.Hour)).
			WithExistingMember(member).
			Build(db)
		sessionIdentifiers[i] = ssn.ID
	}

	testCases := []struct {
		name            string
		size            int32
		token           string
		expectedCount   int
		expectedFirstID string
	}{
		{
			name:          "size of 0 returns all items",
			size:          0,
			expectedCount: SessionsToMake,
		},
		{
			name:          "size of < 0 returns all items",
			size:          -1,
			expectedCount: SessionsToMake,
		},
		{
			name:          "size over 0 returns n items",
			size:          4,
			expectedCount: 4,
		},
		{
			name:          "size of 1 returns single item",
			size:          1,
			expectedCount: 1,
		},
		{
			name:          "size over total items returns all items",
			size:          51,
			expectedCount: SessionsToMake,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := sessionServer.ListSessionsForUser(context.Background(), &sessionpb.ListSessionsForUserRequest{
				UserId:    member.ID,
				PageSize:  tc.size,
				PageToken: tc.token,
			})
			require.NoError(t, err)
			require.NotNil(t, res)
			require.Len(t, res.Sessions, tc.expectedCount)
			require.Len(t, sessionIdentifiers, 50)

			if tc.expectedFirstID != "" {
				assert.Equal(t, tc.expectedFirstID, res.Sessions[0].SessionId)
			}
		})
	}
}

func TestCreateSession_Success(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })
	sessionRepo := session.New(db)

	fakeUserID := uuid.NewString()
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	fakeUserClient := fake.NewFakeUserClient().WithUser(fakeUserID, "user", "username")
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

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

	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	fakeUserClient := fake.NewFakeUserClient().WithUser(uuid.NewString(), "user", "username")
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

	req := &sessionpb.CreateSessionRequest{
		UserId: uuid.NewString(),
		Name:   "New Session",
	}

	resp, err := sessionServer.CreateSession(context.Background(), req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, server.ErrUserNotFound)
	assert.Nil(t, resp)
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
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

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

	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	fakeUserClient := fake.NewFakeUserClient().WithUser(testUserID, "user", "username")
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

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
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	fakeUserClient := fake.NewFakeUserClient().WithUser(testUserID, "user", "username")
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

	req := &sessionpb.AddMemberToSessionRequest{
		SessionId: uuid.NewString(),
		UserId:    testUserID,
	}

	_, err := sessionServer.AddMemberToSession(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, server.ErrSessionNotFound)
}

func TestAddMemberToSession_WhenUserNotFound_Errors(t *testing.T) {
	db := testinfra.SetupTestDB(t, "../db/migrations")
	t.Cleanup(func() { db.Close() })
	sessionRepo := session.New(db)

	ssn := builder.NewSessionBuilder(t).
		WithName("Test Session").
		Build(db)

	testUserID := uuid.NewString()
	fakeUserClient := fake.NewFakeUserClient()
	fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
	sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

	req := &sessionpb.AddMemberToSessionRequest{
		SessionId: ssn.ID,
		UserId:    testUserID,
	}

	_, err := sessionServer.AddMemberToSession(context.Background(), req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, server.ErrUserNotFound)
}

func TestUpdateSessionMemberAdminState(t *testing.T) {
	ownerID := uuid.NewString()
	memberID := uuid.NewString()

	testCases := []struct {
		name           string
		memberToUpdate string
		value          bool
		expectedAdmin  bool
		expectedErr    error
	}{
		{
			name:           "normal case, updating member to admin",
			memberToUpdate: memberID,
			value:          true,
			expectedAdmin:  true,
		},
		{
			name:           "removing fails when session only has one admin",
			memberToUpdate: ownerID,
			value:          false,
			expectedAdmin:  true,
			expectedErr:    server.ErrSessionMustHaveAtLeastOneAdmin,
		},
	}

	q := "select is_admin from session_members where member_id = ?;"
	for _, tc := range testCases {
		db := testinfra.SetupTestDB(t, "../db/migrations")

		sessionRepo := session.New(db)
		fakeUserClient := fake.NewFakeUserClient()
		fakePublisher := fake.NewFakeSessionMemberAddedPublisher()
		sessionServer := server.NewSessionServer(db, sessionRepo, fakeUserClient, fakePublisher, slog.Default())

		ssn := builder.NewSessionBuilder(t).
			WithID(uuid.New()).
			WithName("Test Session").
			WithMember(builder.SessionMemberParams{
				ID:       ownerID,
				Name:     "owner",
				Username: "owner",
				IsOwner:  true,
				IsAdmin:  true,
			}).
			WithMember(builder.SessionMemberParams{
				ID:       memberID,
				Name:     "member",
				Username: "member",
			}).
			Build(db)

		t.Run(tc.name, func(t *testing.T) {
			_, err := sessionServer.UpdateSessionMemberAdminState(context.Background(), &sessionpb.UpdateSessionMemberAdminStateRequest{
				SessionId: ssn.ID,
				UserId:    tc.memberToUpdate,
				IsAdmin:   tc.value,
			})

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, tc.expectedErr, err)
			}

			var isAdmin bool
			err = db.QueryRow(q, tc.memberToUpdate).Scan(&isAdmin)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedAdmin, isAdmin)
		})

		db.Close()
	}
}

func assertUserInsertedAsMember(t *testing.T, db *sql.DB, sessionID, memberID, expectedName, expectedUsername string, expectedOwner bool) {
	t.Helper()

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

func assertStatusHasCode(t *testing.T, err error, expectedCode codes.Code) {
	t.Helper()

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, expectedCode, st.Code())
}
