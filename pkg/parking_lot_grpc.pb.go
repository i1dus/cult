// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.3
// source: parking_lot.proto

package parking_lot

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
	ParkingAPI_GetParkingLot_FullMethodName   = "/api.ParkingAPI/GetParkingLot"
	ParkingAPI_ListParkingLots_FullMethodName = "/api.ParkingAPI/ListParkingLots"
)

// ParkingAPIClient is the client API for ParkingAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ParkingAPIClient interface {
	GetParkingLot(ctx context.Context, in *GetParkingLotRequest, opts ...grpc.CallOption) (*GetParkingLotResponse, error)
	ListParkingLots(ctx context.Context, in *ListParkingLotsRequest, opts ...grpc.CallOption) (*ListParkingLotsResponse, error)
}

type parkingAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewParkingAPIClient(cc grpc.ClientConnInterface) ParkingAPIClient {
	return &parkingAPIClient{cc}
}

func (c *parkingAPIClient) GetParkingLot(ctx context.Context, in *GetParkingLotRequest, opts ...grpc.CallOption) (*GetParkingLotResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetParkingLotResponse)
	err := c.cc.Invoke(ctx, ParkingAPI_GetParkingLot_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *parkingAPIClient) ListParkingLots(ctx context.Context, in *ListParkingLotsRequest, opts ...grpc.CallOption) (*ListParkingLotsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListParkingLotsResponse)
	err := c.cc.Invoke(ctx, ParkingAPI_ListParkingLots_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ParkingAPIServer is the server API for ParkingAPI service.
// All implementations must embed UnimplementedParkingAPIServer
// for forward compatibility.
type ParkingAPIServer interface {
	GetParkingLot(context.Context, *GetParkingLotRequest) (*GetParkingLotResponse, error)
	ListParkingLots(context.Context, *ListParkingLotsRequest) (*ListParkingLotsResponse, error)
	mustEmbedUnimplementedParkingAPIServer()
}

// UnimplementedParkingAPIServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedParkingAPIServer struct{}

func (UnimplementedParkingAPIServer) GetParkingLot(context.Context, *GetParkingLotRequest) (*GetParkingLotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetParkingLot not implemented")
}
func (UnimplementedParkingAPIServer) ListParkingLots(context.Context, *ListParkingLotsRequest) (*ListParkingLotsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListParkingLots not implemented")
}
func (UnimplementedParkingAPIServer) mustEmbedUnimplementedParkingAPIServer() {}
func (UnimplementedParkingAPIServer) testEmbeddedByValue()                    {}

// UnsafeParkingAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ParkingAPIServer will
// result in compilation errors.
type UnsafeParkingAPIServer interface {
	mustEmbedUnimplementedParkingAPIServer()
}

func RegisterParkingAPIServer(s grpc.ServiceRegistrar, srv ParkingAPIServer) {
	// If the following call pancis, it indicates UnimplementedParkingAPIServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ParkingAPI_ServiceDesc, srv)
}

func _ParkingAPI_GetParkingLot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetParkingLotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParkingAPIServer).GetParkingLot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ParkingAPI_GetParkingLot_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParkingAPIServer).GetParkingLot(ctx, req.(*GetParkingLotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ParkingAPI_ListParkingLots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListParkingLotsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ParkingAPIServer).ListParkingLots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ParkingAPI_ListParkingLots_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ParkingAPIServer).ListParkingLots(ctx, req.(*ListParkingLotsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ParkingAPI_ServiceDesc is the grpc.ServiceDesc for ParkingAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ParkingAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.ParkingAPI",
	HandlerType: (*ParkingAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetParkingLot",
			Handler:    _ParkingAPI_GetParkingLot_Handler,
		},
		{
			MethodName: "ListParkingLots",
			Handler:    _ParkingAPI_ListParkingLots_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "parking_lot.proto",
}
