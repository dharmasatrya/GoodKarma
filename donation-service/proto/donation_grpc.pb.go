// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: donation.proto

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
	DonationService_CreateDonation_FullMethodName             = "/donation.DonationService/CreateDonation"
	DonationService_UpdateDonationStatus_FullMethodName       = "/donation.DonationService/UpdateDonationStatus"
	DonationService_UpdateDonationStatusXendit_FullMethodName = "/donation.DonationService/UpdateDonationStatusXendit"
	DonationService_GetDonationsByUserId_FullMethodName       = "/donation.DonationService/GetDonationsByUserId"
	DonationService_GetDonationsByEventId_FullMethodName      = "/donation.DonationService/GetDonationsByEventId"
)

// DonationServiceClient is the client API for DonationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DonationServiceClient interface {
	CreateDonation(ctx context.Context, in *CreateDonationRequest, opts ...grpc.CallOption) (*CreateDonationResponse, error)
	UpdateDonationStatus(ctx context.Context, in *UpdateDonationStatusRequest, opts ...grpc.CallOption) (*UpdateDonationStatusResponse, error)
	UpdateDonationStatusXendit(ctx context.Context, in *UpdateDonationStatusRequest, opts ...grpc.CallOption) (*UpdateDonationStatusResponse, error)
	GetDonationsByUserId(ctx context.Context, in *GetDonationsByUserIdRequest, opts ...grpc.CallOption) (*GetDonationsByUserIdResponse, error)
	GetDonationsByEventId(ctx context.Context, in *GetDonationsByEventIdRequest, opts ...grpc.CallOption) (*GetDonationsByEventIdResponse, error)
}

type donationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDonationServiceClient(cc grpc.ClientConnInterface) DonationServiceClient {
	return &donationServiceClient{cc}
}

func (c *donationServiceClient) CreateDonation(ctx context.Context, in *CreateDonationRequest, opts ...grpc.CallOption) (*CreateDonationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateDonationResponse)
	err := c.cc.Invoke(ctx, DonationService_CreateDonation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *donationServiceClient) UpdateDonationStatus(ctx context.Context, in *UpdateDonationStatusRequest, opts ...grpc.CallOption) (*UpdateDonationStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateDonationStatusResponse)
	err := c.cc.Invoke(ctx, DonationService_UpdateDonationStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *donationServiceClient) UpdateDonationStatusXendit(ctx context.Context, in *UpdateDonationStatusRequest, opts ...grpc.CallOption) (*UpdateDonationStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateDonationStatusResponse)
	err := c.cc.Invoke(ctx, DonationService_UpdateDonationStatusXendit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *donationServiceClient) GetDonationsByUserId(ctx context.Context, in *GetDonationsByUserIdRequest, opts ...grpc.CallOption) (*GetDonationsByUserIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDonationsByUserIdResponse)
	err := c.cc.Invoke(ctx, DonationService_GetDonationsByUserId_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *donationServiceClient) GetDonationsByEventId(ctx context.Context, in *GetDonationsByEventIdRequest, opts ...grpc.CallOption) (*GetDonationsByEventIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDonationsByEventIdResponse)
	err := c.cc.Invoke(ctx, DonationService_GetDonationsByEventId_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DonationServiceServer is the server API for DonationService service.
// All implementations must embed UnimplementedDonationServiceServer
// for forward compatibility.
type DonationServiceServer interface {
	CreateDonation(context.Context, *CreateDonationRequest) (*CreateDonationResponse, error)
	UpdateDonationStatus(context.Context, *UpdateDonationStatusRequest) (*UpdateDonationStatusResponse, error)
	UpdateDonationStatusXendit(context.Context, *UpdateDonationStatusRequest) (*UpdateDonationStatusResponse, error)
	GetDonationsByUserId(context.Context, *GetDonationsByUserIdRequest) (*GetDonationsByUserIdResponse, error)
	GetDonationsByEventId(context.Context, *GetDonationsByEventIdRequest) (*GetDonationsByEventIdResponse, error)
	mustEmbedUnimplementedDonationServiceServer()
}

// UnimplementedDonationServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDonationServiceServer struct{}

func (UnimplementedDonationServiceServer) CreateDonation(context.Context, *CreateDonationRequest) (*CreateDonationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDonation not implemented")
}
func (UnimplementedDonationServiceServer) UpdateDonationStatus(context.Context, *UpdateDonationStatusRequest) (*UpdateDonationStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDonationStatus not implemented")
}
func (UnimplementedDonationServiceServer) UpdateDonationStatusXendit(context.Context, *UpdateDonationStatusRequest) (*UpdateDonationStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDonationStatusXendit not implemented")
}
func (UnimplementedDonationServiceServer) GetDonationsByUserId(context.Context, *GetDonationsByUserIdRequest) (*GetDonationsByUserIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDonationsByUserId not implemented")
}
func (UnimplementedDonationServiceServer) GetDonationsByEventId(context.Context, *GetDonationsByEventIdRequest) (*GetDonationsByEventIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDonationsByEventId not implemented")
}
func (UnimplementedDonationServiceServer) mustEmbedUnimplementedDonationServiceServer() {}
func (UnimplementedDonationServiceServer) testEmbeddedByValue()                         {}

// UnsafeDonationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DonationServiceServer will
// result in compilation errors.
type UnsafeDonationServiceServer interface {
	mustEmbedUnimplementedDonationServiceServer()
}

func RegisterDonationServiceServer(s grpc.ServiceRegistrar, srv DonationServiceServer) {
	// If the following call pancis, it indicates UnimplementedDonationServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DonationService_ServiceDesc, srv)
}

func _DonationService_CreateDonation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDonationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DonationServiceServer).CreateDonation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DonationService_CreateDonation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DonationServiceServer).CreateDonation(ctx, req.(*CreateDonationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DonationService_UpdateDonationStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDonationStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DonationServiceServer).UpdateDonationStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DonationService_UpdateDonationStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DonationServiceServer).UpdateDonationStatus(ctx, req.(*UpdateDonationStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DonationService_UpdateDonationStatusXendit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDonationStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DonationServiceServer).UpdateDonationStatusXendit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DonationService_UpdateDonationStatusXendit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DonationServiceServer).UpdateDonationStatusXendit(ctx, req.(*UpdateDonationStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DonationService_GetDonationsByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDonationsByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DonationServiceServer).GetDonationsByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DonationService_GetDonationsByUserId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DonationServiceServer).GetDonationsByUserId(ctx, req.(*GetDonationsByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DonationService_GetDonationsByEventId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDonationsByEventIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DonationServiceServer).GetDonationsByEventId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DonationService_GetDonationsByEventId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DonationServiceServer).GetDonationsByEventId(ctx, req.(*GetDonationsByEventIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DonationService_ServiceDesc is the grpc.ServiceDesc for DonationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DonationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "donation.DonationService",
	HandlerType: (*DonationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateDonation",
			Handler:    _DonationService_CreateDonation_Handler,
		},
		{
			MethodName: "UpdateDonationStatus",
			Handler:    _DonationService_UpdateDonationStatus_Handler,
		},
		{
			MethodName: "UpdateDonationStatusXendit",
			Handler:    _DonationService_UpdateDonationStatusXendit_Handler,
		},
		{
			MethodName: "GetDonationsByUserId",
			Handler:    _DonationService_GetDonationsByUserId_Handler,
		},
		{
			MethodName: "GetDonationsByEventId",
			Handler:    _DonationService_GetDonationsByEventId_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "donation.proto",
}
