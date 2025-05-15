package errors

import "errors"

var (
	ErrSessionNotFound                 = errors.New("session not found")
	ErrCannotUpdateInactiveSession     = errors.New("cannot update inactive session")
	ErrSessionMustHaveAtLeastOneMember = errors.New("session must have at least one member")
	ErrSessionMustHaveAtLeastOneAdmin  = errors.New("session must have at least one admin")
	ErrSessionMemberNotFound           = errors.New("session member not found")
)
