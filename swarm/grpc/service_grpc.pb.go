// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.0
// source: swarm/grpc/service.proto

package grpc

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

// WorkerServiceClient is the client API for WorkerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WorkerServiceClient interface {
	Init(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*InitReply, error)
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetTorrents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TorrentsReply, error)
	GetTorrentScore(ctx context.Context, in *TorrentScoreRequest, opts ...grpc.CallOption) (*TorrentScoreReply, error)
	DropTorrent(ctx context.Context, in *TorrentDropRequest, opts ...grpc.CallOption) (*TorrentDropReply, error)
	// TODO
	UpdateTorrent(ctx context.Context, in *TorrentUpdateRequest, opts ...grpc.CallOption) (*TorrentUpdateReply, error)
	SaveTorrentFile(ctx context.Context, in *TFileSaveRequest, opts ...grpc.CallOption) (*TFileSaveReply, error)
	GetSystemFreeSpace(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SystemSpaceReply, error)
}

type workerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkerServiceClient(cc grpc.ClientConnInterface) WorkerServiceClient {
	return &workerServiceClient{cc}
}

func (c *workerServiceClient) Init(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*InitReply, error) {
	out := new(InitReply)
	err := c.cc.Invoke(ctx, "/grpc.WorkerService/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerServiceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/grpc.WorkerService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerServiceClient) GetTorrents(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TorrentsReply, error) {
	out := new(TorrentsReply)
	err := c.cc.Invoke(ctx, "/grpc.WorkerService/GetTorrents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerServiceClient) GetTorrentScore(ctx context.Context, in *TorrentScoreRequest, opts ...grpc.CallOption) (*TorrentScoreReply, error) {
	out := new(TorrentScoreReply)
	err := c.cc.Invoke(ctx, "/grpc.WorkerService/GetTorrentScore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerServiceClient) DropTorrent(ctx context.Context, in *TorrentDropRequest, opts ...grpc.CallOption) (*TorrentDropReply, error) {
	out := new(TorrentDropReply)
	err := c.cc.Invoke(ctx, "/grpc.WorkerService/DropTorrent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerServiceClient) UpdateTorrent(ctx context.Context, in *TorrentUpdateRequest, opts ...grpc.CallOption) (*TorrentUpdateReply, error) {
	out := new(TorrentUpdateReply)
	err := c.cc.Invoke(ctx, "/grpc.WorkerService/UpdateTorrent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerServiceClient) SaveTorrentFile(ctx context.Context, in *TFileSaveRequest, opts ...grpc.CallOption) (*TFileSaveReply, error) {
	out := new(TFileSaveReply)
	err := c.cc.Invoke(ctx, "/grpc.WorkerService/SaveTorrentFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workerServiceClient) GetSystemFreeSpace(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SystemSpaceReply, error) {
	out := new(SystemSpaceReply)
	err := c.cc.Invoke(ctx, "/grpc.WorkerService/GetSystemFreeSpace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorkerServiceServer is the server API for WorkerService service.
// All implementations must embed UnimplementedWorkerServiceServer
// for forward compatibility
type WorkerServiceServer interface {
	Init(context.Context, *emptypb.Empty) (*InitReply, error)
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	GetTorrents(context.Context, *emptypb.Empty) (*TorrentsReply, error)
	GetTorrentScore(context.Context, *TorrentScoreRequest) (*TorrentScoreReply, error)
	DropTorrent(context.Context, *TorrentDropRequest) (*TorrentDropReply, error)
	// TODO
	UpdateTorrent(context.Context, *TorrentUpdateRequest) (*TorrentUpdateReply, error)
	SaveTorrentFile(context.Context, *TFileSaveRequest) (*TFileSaveReply, error)
	GetSystemFreeSpace(context.Context, *emptypb.Empty) (*SystemSpaceReply, error)
	mustEmbedUnimplementedWorkerServiceServer()
}

// UnimplementedWorkerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWorkerServiceServer struct {
}

func (UnimplementedWorkerServiceServer) Init(context.Context, *emptypb.Empty) (*InitReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (UnimplementedWorkerServiceServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedWorkerServiceServer) GetTorrents(context.Context, *emptypb.Empty) (*TorrentsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTorrents not implemented")
}
func (UnimplementedWorkerServiceServer) GetTorrentScore(context.Context, *TorrentScoreRequest) (*TorrentScoreReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTorrentScore not implemented")
}
func (UnimplementedWorkerServiceServer) DropTorrent(context.Context, *TorrentDropRequest) (*TorrentDropReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DropTorrent not implemented")
}
func (UnimplementedWorkerServiceServer) UpdateTorrent(context.Context, *TorrentUpdateRequest) (*TorrentUpdateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTorrent not implemented")
}
func (UnimplementedWorkerServiceServer) SaveTorrentFile(context.Context, *TFileSaveRequest) (*TFileSaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveTorrentFile not implemented")
}
func (UnimplementedWorkerServiceServer) GetSystemFreeSpace(context.Context, *emptypb.Empty) (*SystemSpaceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSystemFreeSpace not implemented")
}
func (UnimplementedWorkerServiceServer) mustEmbedUnimplementedWorkerServiceServer() {}

// UnsafeWorkerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WorkerServiceServer will
// result in compilation errors.
type UnsafeWorkerServiceServer interface {
	mustEmbedUnimplementedWorkerServiceServer()
}

func RegisterWorkerServiceServer(s grpc.ServiceRegistrar, srv WorkerServiceServer) {
	s.RegisterService(&WorkerService_ServiceDesc, srv)
}

func _WorkerService_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.WorkerService/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).Init(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.WorkerService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerService_GetTorrents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).GetTorrents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.WorkerService/GetTorrents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).GetTorrents(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerService_GetTorrentScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TorrentScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).GetTorrentScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.WorkerService/GetTorrentScore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).GetTorrentScore(ctx, req.(*TorrentScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerService_DropTorrent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TorrentDropRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).DropTorrent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.WorkerService/DropTorrent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).DropTorrent(ctx, req.(*TorrentDropRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerService_UpdateTorrent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TorrentUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).UpdateTorrent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.WorkerService/UpdateTorrent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).UpdateTorrent(ctx, req.(*TorrentUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerService_SaveTorrentFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TFileSaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).SaveTorrentFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.WorkerService/SaveTorrentFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).SaveTorrentFile(ctx, req.(*TFileSaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkerService_GetSystemFreeSpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkerServiceServer).GetSystemFreeSpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.WorkerService/GetSystemFreeSpace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkerServiceServer).GetSystemFreeSpace(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// WorkerService_ServiceDesc is the grpc.ServiceDesc for WorkerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WorkerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.WorkerService",
	HandlerType: (*WorkerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _WorkerService_Init_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _WorkerService_Ping_Handler,
		},
		{
			MethodName: "GetTorrents",
			Handler:    _WorkerService_GetTorrents_Handler,
		},
		{
			MethodName: "GetTorrentScore",
			Handler:    _WorkerService_GetTorrentScore_Handler,
		},
		{
			MethodName: "DropTorrent",
			Handler:    _WorkerService_DropTorrent_Handler,
		},
		{
			MethodName: "UpdateTorrent",
			Handler:    _WorkerService_UpdateTorrent_Handler,
		},
		{
			MethodName: "SaveTorrentFile",
			Handler:    _WorkerService_SaveTorrentFile_Handler,
		},
		{
			MethodName: "GetSystemFreeSpace",
			Handler:    _WorkerService_GetSystemFreeSpace_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "swarm/grpc/service.proto",
}
