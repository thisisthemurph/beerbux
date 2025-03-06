package server_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/transaction-service/internal/server"
	"github.com/thisisthemurph/beerbux/transaction-service/protos/transactionpb"
	"github.com/thisisthemurph/beerbux/transaction-service/tests/fake"
)

func TestCreateTransaction_Success(t *testing.T) {
	sessionID := uuid.NewString()
	member1ID := uuid.NewString()
	member2ID := uuid.NewString()

	fakeSessionClient := fake.NewFakeSessionClientBuilder().
		WithSession(&sessionpb.SessionResponse{
			SessionId: sessionID,
			Name:      "Session Name",
			IsActive:  true,
			Members: []*sessionpb.SessionMember{{
				UserId:   member1ID,
				Name:     "Member 1",
				Username: "member1",
			}, {
				UserId:   member2ID,
				Name:     "Member 2",
				Username: "member2",
			}},
		}).
		Build()

	fakeTransactionCreatedPublisher := fake.NewFakeTransactionCreatedPublisher()
	srv := server.NewTransactionServer(fakeSessionClient, fakeTransactionCreatedPublisher)

	res, err := srv.CreateTransaction(context.Background(), &transactionpb.CreateTransactionRequest{
		CreatorId: member1ID,
		SessionId: sessionID,
		MemberAmounts: []*transactionpb.MemberAmount{
			{UserId: member2ID, Amount: 1},
		},
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)

	_, err = uuid.Parse(res.TransactionId)
	assert.NoError(t, err)

	assert.Equal(t, member1ID, res.CreatorId)
	assert.Equal(t, sessionID, res.SessionId)
}

func TestCreateTransaction_WithInvalidInputs_ReturnsError(t *testing.T) {
	sessionID := uuid.NewString()
	member1ID := uuid.NewString()
	member2ID := uuid.NewString()

	testCases := []struct {
		name                string
		sessionID           string
		creatorID           string
		transactionMemberID string
		inactiveSession     bool
		amount              float64
		expectedErr         error
	}{
		{
			name:                "empty session id",
			sessionID:           "",
			creatorID:           member1ID,
			transactionMemberID: member2ID,
			amount:              1,
			expectedErr:         server.ErrSessionIDRequired,
		},
		{
			name:                "invalid session id",
			sessionID:           "session-id",
			creatorID:           member1ID,
			transactionMemberID: member2ID,
			amount:              1,
			expectedErr:         server.ErrInvalidUUID,
		},
		{
			name:                "missing member amount id",
			sessionID:           sessionID,
			creatorID:           member1ID,
			transactionMemberID: "",
			amount:              1,
			expectedErr:         server.ErrMemberAmountUserIDRequired,
		},
		{
			name:                "invalid member amount id",
			sessionID:           sessionID,
			creatorID:           member1ID,
			transactionMemberID: "member-id",
			amount:              1,
			expectedErr:         server.ErrInvalidUUID,
		},
		{
			name:                "member amount too low 1",
			sessionID:           sessionID,
			creatorID:           member1ID,
			transactionMemberID: member2ID,
			amount:              0,
			expectedErr:         server.ErrMemberAmountTooLow,
		},
		{
			name:                "member amount too low 2",
			sessionID:           sessionID,
			creatorID:           member1ID,
			transactionMemberID: member2ID,
			amount:              0.5,
			expectedErr:         server.ErrMemberAmountTooLow,
		},
		{
			name:                "member amount too high 1",
			sessionID:           sessionID,
			creatorID:           member1ID,
			transactionMemberID: member2ID,
			amount:              2.5,
			expectedErr:         server.ErrMemberAmountTooHigh,
		},
		{
			name:                "member amount too high 2",
			sessionID:           sessionID,
			creatorID:           member1ID,
			transactionMemberID: member2ID,
			amount:              20,
			expectedErr:         server.ErrMemberAmountTooHigh,
		},
		{
			name:                "creator is in member amounts",
			sessionID:           sessionID,
			creatorID:           member1ID,
			transactionMemberID: member1ID,
			amount:              1,
			expectedErr:         server.ErrCreatorCannotBeMember,
		},
		{
			name:                "inactive session",
			sessionID:           sessionID,
			creatorID:           member1ID,
			transactionMemberID: member2ID,
			amount:              1,
			inactiveSession:     true,
			expectedErr:         server.ErrInactiveSession,
		},
		{
			name:                "member is not member of session",
			sessionID:           sessionID,
			creatorID:           member1ID,
			transactionMemberID: uuid.NewString(),
			amount:              1,
			expectedErr:         server.ErrMemberNotPartOfSession,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeSessionClient := fake.NewFakeSessionClientBuilder().
				WithSession(&sessionpb.SessionResponse{
					SessionId: sessionID,
					Name:      "Session Name",
					IsActive:  !tc.inactiveSession,
					Members: []*sessionpb.SessionMember{{
						UserId:   member1ID,
						Name:     "Member 1",
						Username: "member1",
					}, {
						UserId:   member2ID,
						Name:     "Member 2",
						Username: "member2",
					}},
				}).
				Build()

			fakeTransactionCreatedPublisher := fake.NewFakeTransactionCreatedPublisher()
			srv := server.NewTransactionServer(fakeSessionClient, fakeTransactionCreatedPublisher)

			_, err := srv.CreateTransaction(context.Background(), &transactionpb.CreateTransactionRequest{
				CreatorId: tc.creatorID,
				SessionId: tc.sessionID,
				MemberAmounts: []*transactionpb.MemberAmount{
					{UserId: tc.transactionMemberID, Amount: tc.amount},
				},
			})

			assert.Error(t, err)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestCreateTransaction_WithNoMemberAmounts_Failure(t *testing.T) {
	sessionID := uuid.NewString()
	member1ID := uuid.NewString()
	member2ID := uuid.NewString()

	fakeSessionClient := fake.NewFakeSessionClientBuilder().
		WithSession(&sessionpb.SessionResponse{
			SessionId: sessionID,
			Name:      "Session Name",
			IsActive:  true,
			Members: []*sessionpb.SessionMember{{
				UserId:   member1ID,
				Name:     "Member 1",
				Username: "member1",
			}, {
				UserId:   member2ID,
				Name:     "Member 2",
				Username: "member2",
			}},
		}).
		Build()

	fakeTransactionCreatedPublisher := fake.NewFakeTransactionCreatedPublisher()
	srv := server.NewTransactionServer(fakeSessionClient, fakeTransactionCreatedPublisher)

	testCases := []struct {
		name          string
		memberAmounts []*transactionpb.MemberAmount
	}{
		{
			name:          "nil member amounts",
			memberAmounts: nil,
		},
		{
			name:          "empty member amounts",
			memberAmounts: []*transactionpb.MemberAmount{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := srv.CreateTransaction(context.Background(), &transactionpb.CreateTransactionRequest{
				CreatorId:     member1ID,
				SessionId:     sessionID,
				MemberAmounts: tc.memberAmounts,
			})

			assert.Error(t, err)
			assert.ErrorIs(t, err, server.ErrMemberAmountRequired)
		})
	}
}
