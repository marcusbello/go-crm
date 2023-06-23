// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: crm.proto

package go_crm

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CRMClient is the client API for CRM service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CRMClient interface {
	// Adds customers to the CRM.
	AddCustomers(ctx context.Context, in *AddCustomersReq, opts ...grpc.CallOption) (*AddCustomersResp, error)
	// Updates customers entries in the CRM.
	UpdateCustomers(ctx context.Context, in *UpdateCustomersReq, opts ...grpc.CallOption) (*UpdateCustomersResp, error)
	// Deletes customers from the CRM.
	DeleteCustomers(ctx context.Context, in *DeleteCustomersReq, opts ...grpc.CallOption) (*DeleteCustomersResp, error)
	// Finds customers in the CRM.
	SearchCustomers(ctx context.Context, in *SearchCustomersReq, opts ...grpc.CallOption) (CRM_SearchCustomersClient, error)
}

type cRMClient struct {
	cc grpc.ClientConnInterface
}

func NewCRMClient(cc grpc.ClientConnInterface) CRMClient {
	return &cRMClient{cc}
}

func (c *cRMClient) AddCustomers(ctx context.Context, in *AddCustomersReq, opts ...grpc.CallOption) (*AddCustomersResp, error) {
	out := new(AddCustomersResp)
	err := c.cc.Invoke(ctx, "/crm.CRM/AddCustomers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cRMClient) UpdateCustomers(ctx context.Context, in *UpdateCustomersReq, opts ...grpc.CallOption) (*UpdateCustomersResp, error) {
	out := new(UpdateCustomersResp)
	err := c.cc.Invoke(ctx, "/crm.CRM/UpdateCustomers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cRMClient) DeleteCustomers(ctx context.Context, in *DeleteCustomersReq, opts ...grpc.CallOption) (*DeleteCustomersResp, error) {
	out := new(DeleteCustomersResp)
	err := c.cc.Invoke(ctx, "/crm.CRM/DeleteCustomers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cRMClient) SearchCustomers(ctx context.Context, in *SearchCustomersReq, opts ...grpc.CallOption) (CRM_SearchCustomersClient, error) {
	stream, err := c.cc.NewStream(ctx, &CRM_ServiceDesc.Streams[0], "/crm.CRM/SearchCustomers", opts...)
	if err != nil {
		return nil, err
	}
	x := &cRMSearchCustomersClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CRM_SearchCustomersClient interface {
	Recv() (*Customer, error)
	grpc.ClientStream
}

type cRMSearchCustomersClient struct {
	grpc.ClientStream
}

func (x *cRMSearchCustomersClient) Recv() (*Customer, error) {
	m := new(Customer)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CRMServer is the server API for CRM service.
// All implementations must embed UnimplementedCRMServer
// for forward compatibility
type CRMServer interface {
	// Adds customers to the CRM.
	AddCustomers(context.Context, *AddCustomersReq) (*AddCustomersResp, error)
	// Updates customers entries in the CRM.
	UpdateCustomers(context.Context, *UpdateCustomersReq) (*UpdateCustomersResp, error)
	// Deletes customers from the CRM.
	DeleteCustomers(context.Context, *DeleteCustomersReq) (*DeleteCustomersResp, error)
	// Finds customers in the CRM.
	SearchCustomers(*SearchCustomersReq, CRM_SearchCustomersServer) error
	mustEmbedUnimplementedCRMServer()
}

// UnimplementedCRMServer must be embedded to have forward compatible implementations.
type UnimplementedCRMServer struct {
}

func (UnimplementedCRMServer) AddCustomers(context.Context, *AddCustomersReq) (*AddCustomersResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCustomers not implemented")
}
func (UnimplementedCRMServer) UpdateCustomers(context.Context, *UpdateCustomersReq) (*UpdateCustomersResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCustomers not implemented")
}
func (UnimplementedCRMServer) DeleteCustomers(context.Context, *DeleteCustomersReq) (*DeleteCustomersResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCustomers not implemented")
}
func (UnimplementedCRMServer) SearchCustomers(*SearchCustomersReq, CRM_SearchCustomersServer) error {
	return status.Errorf(codes.Unimplemented, "method SearchCustomers not implemented")
}
func (UnimplementedCRMServer) mustEmbedUnimplementedCRMServer() {}

// UnsafeCRMServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CRMServer will
// result in compilation errors.
type UnsafeCRMServer interface {
	mustEmbedUnimplementedCRMServer()
}

func RegisterCRMServer(s grpc.ServiceRegistrar, srv CRMServer) {
	s.RegisterService(&CRM_ServiceDesc, srv)
}

func _CRM_AddCustomers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCustomersReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CRMServer).AddCustomers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crm.CRM/AddCustomers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CRMServer).AddCustomers(ctx, req.(*AddCustomersReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CRM_UpdateCustomers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCustomersReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CRMServer).UpdateCustomers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crm.CRM/UpdateCustomers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CRMServer).UpdateCustomers(ctx, req.(*UpdateCustomersReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CRM_DeleteCustomers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCustomersReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CRMServer).DeleteCustomers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crm.CRM/DeleteCustomers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CRMServer).DeleteCustomers(ctx, req.(*DeleteCustomersReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CRM_SearchCustomers_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SearchCustomersReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CRMServer).SearchCustomers(m, &cRMSearchCustomersServer{stream})
}

type CRM_SearchCustomersServer interface {
	Send(*Customer) error
	grpc.ServerStream
}

type cRMSearchCustomersServer struct {
	grpc.ServerStream
}

func (x *cRMSearchCustomersServer) Send(m *Customer) error {
	return x.ServerStream.SendMsg(m)
}

// CRM_ServiceDesc is the grpc.ServiceDesc for CRM service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CRM_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "crm.CRM",
	HandlerType: (*CRMServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddCustomers",
			Handler:    _CRM_AddCustomers_Handler,
		},
		{
			MethodName: "UpdateCustomers",
			Handler:    _CRM_UpdateCustomers_Handler,
		},
		{
			MethodName: "DeleteCustomers",
			Handler:    _CRM_DeleteCustomers_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SearchCustomers",
			Handler:       _CRM_SearchCustomers_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "crm.proto",
}