package errors

import "errors"

var (
	ErrSessionNotFound                 = errors.New("session not found")
	ErrCannotUpdateInactiveSession     = errors.New("cannot update inactive session")
	ErrSessionMustHaveAtLeastOneMember = errors.New("session must have at least one member")
	ErrSessionMustHaveAtLeastOneAdmin  = errors.New("session must have at least one admin")
	ErrSessionMemberNotFound           = errors.New("session member not found")
	ErrCreatorCannotBeMember           = errors.New("creator cannot be a member of the transaction")
	ErrInactiveSession                 = errors.New("session is inactive")
	ErrNotAllMembersPartOfSession      = errors.New("member is not part of the session")
	ErrMemberAmountRequired            = errors.New("member amount is required")
	ErrMemberAmountTooLow              = errors.New("member amount must be at least 0.5")
	ErrMemberAmountTooHigh             = errors.New("member amount cannot be more than 2")
)
