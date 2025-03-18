package fake

import (
	"context"
	"errors"
	"github.com/thisisthemurph/beerbux/auth-service/protos/authpb"
	"google.golang.org/grpc"
)

type AuthClient struct {
	response *authpb.RefreshTokenResponse
}

func NewFakeAuthClient(expectedResponse *authpb.RefreshTokenResponse) authpb.AuthClient {
	return &AuthClient{
		response: expectedResponse,
	}
}

func (a AuthClient) Login(ctx context.Context, in *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.LoginResponse, error) {
	panic("implement me")
}

func (a AuthClient) Signup(ctx context.Context, in *authpb.SignupRequest, opts ...grpc.CallOption) (*authpb.SignupResponse, error) {
	panic("implement me")
}

func (a AuthClient) RefreshToken(ctx context.Context, in *authpb.RefreshTokenRequest, opts ...grpc.CallOption) (*authpb.RefreshTokenResponse, error) {
	if a.response == nil {
		return nil, errors.New("no expected response set")
	}
	return a.response, nil
}
