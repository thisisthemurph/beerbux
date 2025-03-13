package fake

import (
	"context"
	"errors"

	"github.com/thisisthemurph/beerbux/user-service/protos/userpb"
	"google.golang.org/grpc"
)

type UserClient struct {
	data []*userpb.GetUserResponse
}

func NewFakeUserClient() *UserClient {
	return &UserClient{}
}

func (c *UserClient) WithUser(id, name, username string) *UserClient {
	usr := &userpb.GetUserResponse{
		UserId:   id,
		Name:     name,
		Username: username,
		Bio:      nil,
	}

	c.data = append(c.data, usr)
	return c
}

func (c *UserClient) GetUser(ctx context.Context, in *userpb.GetUserRequest, opts ...grpc.CallOption) (*userpb.GetUserResponse, error) {
	for _, user := range c.data {
		if user.UserId == in.UserId {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (c *UserClient) CreateUser(ctx context.Context, in *userpb.CreateUserRequest, opts ...grpc.CallOption) (*userpb.UserResponse, error) {
	panic("not implemented")
}

func (c *UserClient) UpdateUser(ctx context.Context, in *userpb.UpdateUserRequest, opts ...grpc.CallOption) (*userpb.UserResponse, error) {
	panic("not implemented")
}
