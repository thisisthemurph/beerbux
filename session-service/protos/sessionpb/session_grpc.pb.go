// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: session.proto

package sessionpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Session_CreateSession_FullMethodName      = "/Session/CreateSession"
	Session_GetSession_FullMethodName         = "/Session/GetSession"
	Session_AddMemberToSession_FullMethodName = "/Session/AddMemberToSession"
)

// SessionClient is the client API for Session service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SessionClient interface {
	CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*SessionResponse, error)
	GetSession(ctx context.Context, in *GetSessionRequest, opts ...grpc.CallOption) (*SessionResponse, error)
	AddMemberToSession(ctx context.Context, in *AddMemberToSessionRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type sessionClient struct {
	cc grpc.ClientConnInterface
}

func NewSessionClient(cc grpc.ClientConnInterface) SessionClient {
	return &sessionClient{cc}
}

func (c *sessionClient) CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*SessionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SessionResponse)
	err := c.cc.Invoke(ctx, Session_CreateSession_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionClient) GetSession(ctx context.Context, in *GetSessionRequest, opts ...grpc.CallOption) (*SessionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SessionResponse)
	err := c.cc.Invoke(ctx, Session_GetSession_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionClient) AddMemberToSession(ctx context.Context, in *AddMemberToSessionRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, Session_AddMemberToSession_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SessionServer is the server API for Session service.
// All implementations must embed UnimplementedSessionServer
// for forward compatibility.
type SessionServer interface {
	CreateSession(context.Context, *CreateSessionRequest) (*SessionResponse, error)
	GetSession(context.Context, *GetSessionRequest) (*SessionResponse, error)
	AddMemberToSession(context.Context, *AddMemberToSessionRequest) (*EmptyResponse, error)
	mustEmbedUnimplementedSessionServer()
}

// UnimplementedSessionServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSessionServer struct{}

func (UnimplementedSessionServer) CreateSession(context.Context, *CreateSessionRequest) (*SessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSession not implemented")
}
func (UnimplementedSessionServer) GetSession(context.Context, *GetSessionRequest) (*SessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSession not implemented")
}
func (UnimplementedSessionServer) AddMemberToSession(context.Context, *AddMemberToSessionRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddMemberToSession not implemented")
}
func (UnimplementedSessionServer) mustEmbedUnimplementedSessionServer() {}
func (UnimplementedSessionServer) testEmbeddedByValue()                 {}

// UnsafeSessionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SessionServer will
// result in compilation errors.
type UnsafeSessionServer interface {
	mustEmbedUnimplementedSessionServer()
}

func RegisterSessionServer(s grpc.ServiceRegistrar, srv SessionServer) {
	// If the following call pancis, it indicates UnimplementedSessionServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Session_ServiceDesc, srv)
}

func _Session_CreateSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionServer).CreateSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Session_CreateSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionServer).CreateSession(ctx, req.(*CreateSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Session_GetSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionServer).GetSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Session_GetSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionServer).GetSession(ctx, req.(*GetSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Session_AddMemberToSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddMemberToSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SessionServer).AddMemberToSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Session_AddMemberToSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SessionServer).AddMemberToSession(ctx, req.(*AddMemberToSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Session_ServiceDesc is the grpc.ServiceDesc for Session service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Session_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Session",
	HandlerType: (*SessionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSession",
			Handler:    _Session_CreateSession_Handler,
		},
		{
			MethodName: "GetSession",
			Handler:    _Session_GetSession_Handler,
		},
		{
			MethodName: "AddMemberToSession",
			Handler:    _Session_AddMemberToSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "session.proto",
}
