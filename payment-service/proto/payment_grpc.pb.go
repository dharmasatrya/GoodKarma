// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: payment.proto

package proto

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
	PaymentService_CreateWallet_FullMethodName        = "/payment.PaymentService/CreateWallet"
	PaymentService_UpdateWalletBalance_FullMethodName = "/payment.PaymentService/UpdateWalletBalance"
	PaymentService_Withdraw_FullMethodName            = "/payment.PaymentService/Withdraw"
	PaymentService_CreateInvoice_FullMethodName       = "/payment.PaymentService/CreateInvoice"
)

// PaymentServiceClient is the client API for PaymentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentServiceClient interface {
	CreateWallet(ctx context.Context, in *CreateWalletRequest, opts ...grpc.CallOption) (*CreateWalletResponse, error)
	UpdateWalletBalance(ctx context.Context, in *UpdateWalletBalanceRequest, opts ...grpc.CallOption) (*UpdateWalleetBalanceResponse, error)
	Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawResponse, error)
	CreateInvoice(ctx context.Context, in *CreateInvoiceRequest, opts ...grpc.CallOption) (*CreateInvoiceResponse, error)
}

type paymentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentServiceClient(cc grpc.ClientConnInterface) PaymentServiceClient {
	return &paymentServiceClient{cc}
}

func (c *paymentServiceClient) CreateWallet(ctx context.Context, in *CreateWalletRequest, opts ...grpc.CallOption) (*CreateWalletResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateWalletResponse)
	err := c.cc.Invoke(ctx, PaymentService_CreateWallet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) UpdateWalletBalance(ctx context.Context, in *UpdateWalletBalanceRequest, opts ...grpc.CallOption) (*UpdateWalleetBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateWalleetBalanceResponse)
	err := c.cc.Invoke(ctx, PaymentService_UpdateWalletBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WithdrawResponse)
	err := c.cc.Invoke(ctx, PaymentService_Withdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) CreateInvoice(ctx context.Context, in *CreateInvoiceRequest, opts ...grpc.CallOption) (*CreateInvoiceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateInvoiceResponse)
	err := c.cc.Invoke(ctx, PaymentService_CreateInvoice_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentServiceServer is the server API for PaymentService service.
// All implementations must embed UnimplementedPaymentServiceServer
// for forward compatibility.
type PaymentServiceServer interface {
	CreateWallet(context.Context, *CreateWalletRequest) (*CreateWalletResponse, error)
	UpdateWalletBalance(context.Context, *UpdateWalletBalanceRequest) (*UpdateWalleetBalanceResponse, error)
	Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error)
	CreateInvoice(context.Context, *CreateInvoiceRequest) (*CreateInvoiceResponse, error)
	mustEmbedUnimplementedPaymentServiceServer()
}

// UnimplementedPaymentServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPaymentServiceServer struct{}

func (UnimplementedPaymentServiceServer) CreateWallet(context.Context, *CreateWalletRequest) (*CreateWalletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWallet not implemented")
}
func (UnimplementedPaymentServiceServer) UpdateWalletBalance(context.Context, *UpdateWalletBalanceRequest) (*UpdateWalleetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateWalletBalance not implemented")
}
func (UnimplementedPaymentServiceServer) Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}
func (UnimplementedPaymentServiceServer) CreateInvoice(context.Context, *CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateInvoice not implemented")
}
func (UnimplementedPaymentServiceServer) mustEmbedUnimplementedPaymentServiceServer() {}
func (UnimplementedPaymentServiceServer) testEmbeddedByValue()                        {}

// UnsafePaymentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentServiceServer will
// result in compilation errors.
type UnsafePaymentServiceServer interface {
	mustEmbedUnimplementedPaymentServiceServer()
}

func RegisterPaymentServiceServer(s grpc.ServiceRegistrar, srv PaymentServiceServer) {
	// If the following call pancis, it indicates UnimplementedPaymentServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PaymentService_ServiceDesc, srv)
}

func _PaymentService_CreateWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).CreateWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_CreateWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).CreateWallet(ctx, req.(*CreateWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_UpdateWalletBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateWalletBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).UpdateWalletBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_UpdateWalletBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).UpdateWalletBalance(ctx, req.(*UpdateWalletBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_Withdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_Withdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).Withdraw(ctx, req.(*WithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PaymentService_CreateInvoice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateInvoiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceServer).CreateInvoice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_CreateInvoice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceServer).CreateInvoice(ctx, req.(*CreateInvoiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PaymentService_ServiceDesc is the grpc.ServiceDesc for PaymentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PaymentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "payment.PaymentService",
	HandlerType: (*PaymentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateWallet",
			Handler:    _PaymentService_CreateWallet_Handler,
		},
		{
			MethodName: "UpdateWalletBalance",
			Handler:    _PaymentService_UpdateWalletBalance_Handler,
		},
		{
			MethodName: "Withdraw",
			Handler:    _PaymentService_Withdraw_Handler,
		},
		{
			MethodName: "CreateInvoice",
			Handler:    _PaymentService_CreateInvoice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payment.proto",
}
