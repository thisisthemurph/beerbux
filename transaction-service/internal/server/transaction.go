package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"github.com/thisisthemurph/beerbux/transaction-service/internal/publisher"
	"github.com/thisisthemurph/beerbux/transaction-service/pkg/fn"
	"github.com/thisisthemurph/beerbux/transaction-service/protos/transactionpb"
)

type TransactionServer struct {
	transactionpb.UnimplementedTransactionServer
	sessionClient               sessionpb.SessionClient
	transactionCreatedPublisher publisher.TransactionCreatedPublisher
}

func NewTransactionServer(sessionClient sessionpb.SessionClient, transactionCreatedPublisher publisher.TransactionCreatedPublisher) *TransactionServer {
	return &TransactionServer{
		sessionClient:               sessionClient,
		transactionCreatedPublisher: transactionCreatedPublisher,
	}
}

func (s *TransactionServer) CreateTransaction(ctx context.Context, r *transactionpb.CreateTransactionRequest) (*transactionpb.TransactionResponse, error) {
	if err := validateCreateTransactionRequest(r); err != nil {
		return nil, err
	}

	memberIDLookup := fn.Map(r.MemberAmounts, func(ma *transactionpb.MemberAmount) string {
		return ma.GetUserId()
	})

	if fn.Contains(memberIDLookup, r.CreatorId) {
		return nil, errors.New("creator cannot be a member of the transaction")
	}

	if err := s.validateSession(ctx, r.SessionId, memberIDLookup); err != nil {
		return nil, err
	}

	transactionID := uuid.NewString()

	err := s.transactionCreatedPublisher.Publish(publisher.TransactionCreatedEventData{
		TransactionID: transactionID,
		CreatorID:     r.CreatorId,
		SessionID:     r.SessionId,
		MemberAmounts: fn.Map(r.MemberAmounts, func(ma *transactionpb.MemberAmount) publisher.TransactionCreatedMemberAmount {
			return publisher.TransactionCreatedMemberAmount{
				UserID: ma.GetUserId(),
				Amount: ma.GetAmount(),
			}
		}),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to publish transaction created event: %w", err)
	}

	return &transactionpb.TransactionResponse{
		TransactionId: transactionID,
		CreatorId:     r.CreatorId,
		SessionId:     r.SessionId,
		CreatedAt:     time.Now().String(),
	}, nil
}

// validateSession validates the session associated with the given sessionID.
//
//   - Validate that the session ID is valid.
//   - Validate that the session exists.
//   - Validate that the session is active.
//   - Validate that the given members are members of the session.
func (s *TransactionServer) validateSession(ctx context.Context, sessionID string, memberIDLookup []string) error {
	if sessionID == "" {
		return errors.New("session_id is required")
	}

	ssn, err := s.sessionClient.GetSession(ctx, &sessionpb.GetSessionRequest{
		SessionId: sessionID,
	})

	if err != nil {
		return err
	}

	if !ssn.IsActive {
		return errors.New("cannot create a session for an inactive session")
	}

	sessionMemberIDs := fn.Map(ssn.Members, func(m *sessionpb.SessionMember) string {
		return m.UserId
	})

	for _, memberID := range memberIDLookup {
		if !fn.Contains(sessionMemberIDs, memberID) {
			return errors.New("member is not part of the session")
		}
	}

	return nil
}

// validateCreateTransactionRequest validates the given request.
func validateCreateTransactionRequest(r *transactionpb.CreateTransactionRequest) error {
	if r.SessionId == "" {
		return errors.New("session_id is required")
	}

	if r.MemberAmounts == nil || len(r.MemberAmounts) == 0 {
		return errors.New("member_amounts is required")
	}

	for _, memberAmount := range r.MemberAmounts {
		if memberAmount.UserId == "" {
			return errors.New("member_amounts.user_id is required")
		}

		if memberAmount.Amount <= 0.5 {
			return errors.New("member_amounts.amount must at least 0.5")
		}

		if memberAmount.Amount > 2 {
			return errors.New("member_amounts.amount cannot be more than 2")
		}
	}

	return nil
}
