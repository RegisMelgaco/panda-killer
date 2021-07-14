// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PandaKillerClient is the client API for PandaKiller service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PandaKillerClient interface {
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
	ListAccounts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAccountListResponse, error)
	GetAccountBalance(ctx context.Context, in *GetAccountBalanceRequest, opts ...grpc.CallOption) (*GetAccountBalanceResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}

type pandaKillerClient struct {
	cc grpc.ClientConnInterface
}

func NewPandaKillerClient(cc grpc.ClientConnInterface) PandaKillerClient {
	return &pandaKillerClient{cc}
}

func (c *pandaKillerClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	out := new(CreateAccountResponse)
	err := c.cc.Invoke(ctx, "/pandakiller.PandaKiller/CreateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pandaKillerClient) ListAccounts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAccountListResponse, error) {
	out := new(GetAccountListResponse)
	err := c.cc.Invoke(ctx, "/pandakiller.PandaKiller/ListAccounts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pandaKillerClient) GetAccountBalance(ctx context.Context, in *GetAccountBalanceRequest, opts ...grpc.CallOption) (*GetAccountBalanceResponse, error) {
	out := new(GetAccountBalanceResponse)
	err := c.cc.Invoke(ctx, "/pandakiller.PandaKiller/GetAccountBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pandaKillerClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/pandakiller.PandaKiller/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PandaKillerServer is the server API for PandaKiller service.
// All implementations should embed UnimplementedPandaKillerServer
// for forward compatibility
type PandaKillerServer interface {
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
	ListAccounts(context.Context, *emptypb.Empty) (*GetAccountListResponse, error)
	GetAccountBalance(context.Context, *GetAccountBalanceRequest) (*GetAccountBalanceResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
}

// UnimplementedPandaKillerServer should be embedded to have forward compatible implementations.
type UnimplementedPandaKillerServer struct {
}

func (UnimplementedPandaKillerServer) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedPandaKillerServer) ListAccounts(context.Context, *emptypb.Empty) (*GetAccountListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccounts not implemented")
}
func (UnimplementedPandaKillerServer) GetAccountBalance(context.Context, *GetAccountBalanceRequest) (*GetAccountBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccountBalance not implemented")
}
func (UnimplementedPandaKillerServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}

// UnsafePandaKillerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PandaKillerServer will
// result in compilation errors.
type UnsafePandaKillerServer interface {
	mustEmbedUnimplementedPandaKillerServer()
}

func RegisterPandaKillerServer(s grpc.ServiceRegistrar, srv PandaKillerServer) {
	s.RegisterService(&PandaKiller_ServiceDesc, srv)
}

func _PandaKiller_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PandaKillerServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pandakiller.PandaKiller/CreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PandaKillerServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PandaKiller_ListAccounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PandaKillerServer).ListAccounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pandakiller.PandaKiller/ListAccounts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PandaKillerServer).ListAccounts(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PandaKiller_GetAccountBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PandaKillerServer).GetAccountBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pandakiller.PandaKiller/GetAccountBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PandaKillerServer).GetAccountBalance(ctx, req.(*GetAccountBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PandaKiller_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PandaKillerServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pandakiller.PandaKiller/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PandaKillerServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PandaKiller_ServiceDesc is the grpc.ServiceDesc for PandaKiller service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PandaKiller_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pandakiller.PandaKiller",
	HandlerType: (*PandaKillerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAccount",
			Handler:    _PandaKiller_CreateAccount_Handler,
		},
		{
			MethodName: "ListAccounts",
			Handler:    _PandaKiller_ListAccounts_Handler,
		},
		{
			MethodName: "GetAccountBalance",
			Handler:    _PandaKiller_GetAccountBalance_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _PandaKiller_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/gateway/rpc/service.proto",
}
