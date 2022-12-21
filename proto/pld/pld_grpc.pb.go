// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: pld/v1/pld.proto

package pld

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

// PldServiceClient is the client API for PldService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PldServiceClient interface {
	CheckInBlacklist(ctx context.Context, in *CheckInBlackListReq, opts ...grpc.CallOption) (*CheckInBlackListRes, error)
}

type pldServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPldServiceClient(cc grpc.ClientConnInterface) PldServiceClient {
	return &pldServiceClient{cc}
}

func (c *pldServiceClient) CheckInBlacklist(ctx context.Context, in *CheckInBlackListReq, opts ...grpc.CallOption) (*CheckInBlackListRes, error) {
	out := new(CheckInBlackListRes)
	err := c.cc.Invoke(ctx, "/pld.PldService/CheckInBlacklist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PldServiceServer is the server API for PldService service.
// All implementations must embed UnimplementedPldServiceServer
// for forward compatibility
type PldServiceServer interface {
	CheckInBlacklist(context.Context, *CheckInBlackListReq) (*CheckInBlackListRes, error)
	mustEmbedUnimplementedPldServiceServer()
}

// UnimplementedPldServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPldServiceServer struct {
}

func (UnimplementedPldServiceServer) CheckInBlacklist(context.Context, *CheckInBlackListReq) (*CheckInBlackListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckInBlacklist not implemented")
}
func (UnimplementedPldServiceServer) mustEmbedUnimplementedPldServiceServer() {}

// UnsafePldServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PldServiceServer will
// result in compilation errors.
type UnsafePldServiceServer interface {
	mustEmbedUnimplementedPldServiceServer()
}

func RegisterPldServiceServer(s grpc.ServiceRegistrar, srv PldServiceServer) {
	s.RegisterService(&PldService_ServiceDesc, srv)
}

func _PldService_CheckInBlacklist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckInBlackListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PldServiceServer).CheckInBlacklist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pld.PldService/CheckInBlacklist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PldServiceServer).CheckInBlacklist(ctx, req.(*CheckInBlackListReq))
	}
	return interceptor(ctx, in, info, handler)
}

// PldService_ServiceDesc is the grpc.ServiceDesc for PldService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PldService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pld.PldService",
	HandlerType: (*PldServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckInBlacklist",
			Handler:    _PldService_CheckInBlacklist_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pld/v1/pld.proto",
}