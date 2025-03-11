package fake

import (
	"context"
	"errors"

	"github.com/thisisthemurph/beerbux/session-service/protos/sessionpb"
	"google.golang.org/grpc"
)

type SessionClientBuilderFake struct {
	sessions []*sessionpb.SessionResponse
}

// NewFakeSessionClientBuilder creates a fake SessionClientBuilderFake.
func NewFakeSessionClientBuilder() *SessionClientBuilderFake {
	return &SessionClientBuilderFake{}
}

// Build returns a fake session client conforming to the sessionpb.SessionClient interface.
func (f *SessionClientBuilderFake) Build() sessionpb.SessionClient {
	return f
}

// WithSession adds a session to the FakeSessionClient.
// This is the data that will be available when using the GetSession function.
func (f *SessionClientBuilderFake) WithSession(ssn *sessionpb.SessionResponse) *SessionClientBuilderFake {
	f.sessions = append(f.sessions, ssn)
	return f
}

func (f *SessionClientBuilderFake) CreateSession(ctx context.Context, in *sessionpb.CreateSessionRequest, opts ...grpc.CallOption) (*sessionpb.SessionResponse, error) {
	panic("implement me")
}

func (f *SessionClientBuilderFake) GetSession(ctx context.Context, in *sessionpb.GetSessionRequest, opts ...grpc.CallOption) (*sessionpb.SessionResponse, error) {
	for _, ssn := range f.sessions {
		if ssn.SessionId == in.SessionId {
			return ssn, nil
		}
	}
	return nil, errors.New("session not found")
}

func (f *SessionClientBuilderFake) AddMemberToSession(ctx context.Context, in *sessionpb.AddMemberToSessionRequest, opts ...grpc.CallOption) (*sessionpb.EmptyResponse, error) {
	panic("implement me")
}
