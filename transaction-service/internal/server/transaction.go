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

var (
	ErrInvalidUUID                = errors.New("invalid UUID")
	ErrCreatorCannotBeMember      = errors.New("creator cannot be a member of the transaction")
	ErrSessionIDRequired          = errors.New("session_id is required")
	ErrInactiveSession            = errors.New("session is inactive")
	ErrMemberNotPartOfSession     = errors.New("member is not part of the session")
	ErrMemberAmountRequired       = errors.New("member amount is required")
	ErrMemberAmountUserIDRequired = errors.New("member amount user ID is required")
	ErrMemberAmountTooLow         = errors.New("member amount must be at least 0.5")
	ErrMemberAmountTooHigh        = errors.New("member amount cannot be more than 2")
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
		return nil, ErrCreatorCannotBeMember
	}

	if err := s.validateSession(ctx, r.SessionId, memberIDLookup); err != nil {
		return nil, err
	}

	transactionID := uuid.NewString()
	err := s.transactionCreatedPublisher.Publish(ctx, publisher.TransactionCreatedEvent{
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
		return ErrSessionIDRequired
	}

	ssn, err := s.sessionClient.GetSession(ctx, &sessionpb.GetSessionRequest{
		SessionId: sessionID,
	})

	if err != nil {
		return err
	}

	if !ssn.IsActive {
		return ErrInactiveSession
	}

	sessionMemberIDs := fn.Map(ssn.Members, func(m *sessionpb.SessionMember) string {
		return m.UserId
	})

	for _, memberID := range memberIDLookup {
		if !fn.Contains(sessionMemberIDs, memberID) {
			return ErrMemberNotPartOfSession
		}
	}

	return nil
}

// validateCreateTransactionRequest validates the given request.
func validateCreateTransactionRequest(r *transactionpb.CreateTransactionRequest) error {
	if r.SessionId == "" {
		return ErrSessionIDRequired
	}

	if err := validateStringUUID(r.SessionId); err != nil {
		return ErrInvalidUUID
	}

	if r.MemberAmounts == nil || len(r.MemberAmounts) == 0 {
		return ErrMemberAmountRequired
	}

	for _, memberAmount := range r.MemberAmounts {
		if memberAmount.UserId == "" {
			return ErrMemberAmountUserIDRequired
		}

		if err := validateStringUUID(memberAmount.UserId); err != nil {
			return ErrInvalidUUID
		}

		if memberAmount.Amount <= 0.5 {
			return ErrMemberAmountTooLow
		}

		if memberAmount.Amount > 2 {
			return ErrMemberAmountTooHigh
		}
	}

	return nil
}
