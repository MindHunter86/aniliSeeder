// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.0
// source: swarm/grpc/service.proto

package grpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RegistrationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// initial phase
	WorkerVersion string `protobuf:"bytes,1,opt,name=WorkerVersion,proto3" json:"WorkerVersion,omitempty"`
	WDFreeSpace   uint64 `protobuf:"varint,2,opt,name=WDFreeSpace,proto3" json:"WDFreeSpace,omitempty"`
	// registration phase
	Torrent []*structpb.Struct `protobuf:"bytes,3,rep,name=Torrent,proto3" json:"Torrent,omitempty"`
}

func (x *RegistrationRequest) Reset() {
	*x = RegistrationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationRequest) ProtoMessage() {}

func (x *RegistrationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationRequest.ProtoReflect.Descriptor instead.
func (*RegistrationRequest) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{0}
}

func (x *RegistrationRequest) GetWorkerVersion() string {
	if x != nil {
		return x.WorkerVersion
	}
	return ""
}

func (x *RegistrationRequest) GetWDFreeSpace() uint64 {
	if x != nil {
		return x.WDFreeSpace
	}
	return 0
}

func (x *RegistrationRequest) GetTorrent() []*structpb.Struct {
	if x != nil {
		return x.Torrent
	}
	return nil
}

type RegistrationReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// setup phase
	Config *structpb.Struct `protobuf:"bytes,2,opt,name=Config,proto3" json:"Config,omitempty"`
}

func (x *RegistrationReply) Reset() {
	*x = RegistrationReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationReply) ProtoMessage() {}

func (x *RegistrationReply) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationReply.ProtoReflect.Descriptor instead.
func (*RegistrationReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{1}
}

func (x *RegistrationReply) GetConfig() *structpb.Struct {
	if x != nil {
		return x.Config
	}
	return nil
}

type TorrentsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Torrent []*structpb.Struct `protobuf:"bytes,1,rep,name=Torrent,proto3" json:"Torrent,omitempty"`
}

func (x *TorrentsReply) Reset() {
	*x = TorrentsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentsReply) ProtoMessage() {}

func (x *TorrentsReply) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TorrentsReply.ProtoReflect.Descriptor instead.
func (*TorrentsReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{2}
}

func (x *TorrentsReply) GetTorrent() []*structpb.Struct {
	if x != nil {
		return x.Torrent
	}
	return nil
}

type TorrentScoreRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash string `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *TorrentScoreRequest) Reset() {
	*x = TorrentScoreRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentScoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentScoreRequest) ProtoMessage() {}

func (x *TorrentScoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TorrentScoreRequest.ProtoReflect.Descriptor instead.
func (*TorrentScoreRequest) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{3}
}

func (x *TorrentScoreRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *TorrentScoreRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type TorrentScoreReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ratio float32 `protobuf:"fixed32,1,opt,name=Ratio,proto3" json:"Ratio,omitempty"`
	Score float32 `protobuf:"fixed32,2,opt,name=Score,proto3" json:"Score,omitempty"`
}

func (x *TorrentScoreReply) Reset() {
	*x = TorrentScoreReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentScoreReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentScoreReply) ProtoMessage() {}

func (x *TorrentScoreReply) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TorrentScoreReply.ProtoReflect.Descriptor instead.
func (*TorrentScoreReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{4}
}

func (x *TorrentScoreReply) GetRatio() float32 {
	if x != nil {
		return x.Ratio
	}
	return 0
}

func (x *TorrentScoreReply) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

type TorrentDropRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash   string `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Name   string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Reason string `protobuf:"bytes,3,opt,name=Reason,proto3" json:"Reason,omitempty"`
}

func (x *TorrentDropRequest) Reset() {
	*x = TorrentDropRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentDropRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentDropRequest) ProtoMessage() {}

func (x *TorrentDropRequest) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TorrentDropRequest.ProtoReflect.Descriptor instead.
func (*TorrentDropRequest) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{5}
}

func (x *TorrentDropRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *TorrentDropRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TorrentDropRequest) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

type TorrentDropReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FreedSpace uint64 `protobuf:"varint,1,opt,name=FreedSpace,proto3" json:"FreedSpace,omitempty"`
	FreeSpace  uint64 `protobuf:"varint,2,opt,name=FreeSpace,proto3" json:"FreeSpace,omitempty"`
}

func (x *TorrentDropReply) Reset() {
	*x = TorrentDropReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentDropReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentDropReply) ProtoMessage() {}

func (x *TorrentDropReply) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TorrentDropReply.ProtoReflect.Descriptor instead.
func (*TorrentDropReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{6}
}

func (x *TorrentDropReply) GetFreedSpace() uint64 {
	if x != nil {
		return x.FreedSpace
	}
	return 0
}

func (x *TorrentDropReply) GetFreeSpace() uint64 {
	if x != nil {
		return x.FreeSpace
	}
	return 0
}

type TFileSaveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *TFileSaveRequest) Reset() {
	*x = TFileSaveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TFileSaveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TFileSaveRequest) ProtoMessage() {}

func (x *TFileSaveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TFileSaveRequest.ProtoReflect.Descriptor instead.
func (*TFileSaveRequest) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{7}
}

func (x *TFileSaveRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type TFileSaveReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilePayload uint64 `protobuf:"varint,1,opt,name=FilePayload,proto3" json:"FilePayload,omitempty"`
}

func (x *TFileSaveReply) Reset() {
	*x = TFileSaveReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TFileSaveReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TFileSaveReply) ProtoMessage() {}

func (x *TFileSaveReply) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TFileSaveReply.ProtoReflect.Descriptor instead.
func (*TFileSaveReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{8}
}

func (x *TFileSaveReply) GetFilePayload() uint64 {
	if x != nil {
		return x.FilePayload
	}
	return 0
}

type SystemSpaceReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FreeSpace uint64 `protobuf:"varint,1,opt,name=FreeSpace,proto3" json:"FreeSpace,omitempty"`
}

func (x *SystemSpaceReply) Reset() {
	*x = SystemSpaceReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SystemSpaceReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SystemSpaceReply) ProtoMessage() {}

func (x *SystemSpaceReply) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SystemSpaceReply.ProtoReflect.Descriptor instead.
func (*SystemSpaceReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{9}
}

func (x *SystemSpaceReply) GetFreeSpace() uint64 {
	if x != nil {
		return x.FreeSpace
	}
	return 0
}

var File_swarm_grpc_service_proto protoreflect.FileDescriptor

var file_swarm_grpc_service_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x77, 0x61, 0x72, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x01, 0x0a, 0x13,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x57, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x57, 0x44, 0x46,
	0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b,
	0x57, 0x44, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x12, 0x31, 0x0a, 0x07, 0x54,
	0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x22, 0x44,
	0x0a, 0x11, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x2f, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x06, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x22, 0x42, 0x0a, 0x0d, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x73,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x31, 0x0a, 0x07, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52,
	0x07, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x22, 0x3d, 0x0a, 0x13, 0x54, 0x6f, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48,
	0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x3f, 0x0a, 0x11, 0x54, 0x6f, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x52, 0x61, 0x74, 0x69, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x52, 0x61, 0x74,
	0x69, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x02, 0x52, 0x05, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x22, 0x54, 0x0a, 0x12, 0x54, 0x6f, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x61,
	0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x50,
	0x0a, 0x10, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x46, 0x72, 0x65, 0x65, 0x64, 0x53, 0x70, 0x61, 0x63, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x46, 0x72, 0x65, 0x65, 0x64, 0x53, 0x70, 0x61,
	0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65,
	0x22, 0x2c, 0x0a, 0x10, 0x54, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x32,
	0x0a, 0x0e, 0x54, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x20, 0x0a, 0x0b, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x22, 0x30, 0x0a, 0x10, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53, 0x70, 0x61, 0x63,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70,
	0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x46, 0x72, 0x65, 0x65, 0x53,
	0x70, 0x61, 0x63, 0x65, 0x32, 0xe2, 0x02, 0x0a, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3c, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x13, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x22, 0x00, 0x12, 0x47, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x19, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54,
	0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x17, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e,
	0x74, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x3f, 0x0a,
	0x0b, 0x44, 0x72, 0x6f, 0x70, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x44, 0x72, 0x6f, 0x70, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x46,
	0x69, 0x6c, 0x65, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x41,
	0x0a, 0x0f, 0x53, 0x61, 0x76, 0x65, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x12, 0x16, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x61,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x54, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22,
	0x00, 0x12, 0x46, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x46, 0x72,
	0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x16, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53, 0x70, 0x61,
	0x63, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x32, 0x8b, 0x01, 0x0a, 0x0d, 0x4d, 0x61,
	0x73, 0x74, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x04, 0x50,
	0x69, 0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x40, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x12, 0x19, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x69, 0x6e, 0x64, 0x48, 0x75, 0x6e, 0x74, 0x65, 0x72,
	0x38, 0x36, 0x2f, 0x61, 0x6e, 0x69, 0x6c, 0x69, 0x53, 0x65, 0x65, 0x64, 0x65, 0x72, 0x2f, 0x73,
	0x77, 0x61, 0x72, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_swarm_grpc_service_proto_rawDescOnce sync.Once
	file_swarm_grpc_service_proto_rawDescData = file_swarm_grpc_service_proto_rawDesc
)

func file_swarm_grpc_service_proto_rawDescGZIP() []byte {
	file_swarm_grpc_service_proto_rawDescOnce.Do(func() {
		file_swarm_grpc_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_swarm_grpc_service_proto_rawDescData)
	})
	return file_swarm_grpc_service_proto_rawDescData
}

var file_swarm_grpc_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_swarm_grpc_service_proto_goTypes = []interface{}{
	(*RegistrationRequest)(nil), // 0: grpc.RegistrationRequest
	(*RegistrationReply)(nil),   // 1: grpc.RegistrationReply
	(*TorrentsReply)(nil),       // 2: grpc.TorrentsReply
	(*TorrentScoreRequest)(nil), // 3: grpc.TorrentScoreRequest
	(*TorrentScoreReply)(nil),   // 4: grpc.TorrentScoreReply
	(*TorrentDropRequest)(nil),  // 5: grpc.TorrentDropRequest
	(*TorrentDropReply)(nil),    // 6: grpc.TorrentDropReply
	(*TFileSaveRequest)(nil),    // 7: grpc.TFileSaveRequest
	(*TFileSaveReply)(nil),      // 8: grpc.TFileSaveReply
	(*SystemSpaceReply)(nil),    // 9: grpc.SystemSpaceReply
	(*structpb.Struct)(nil),     // 10: google.protobuf.Struct
	(*emptypb.Empty)(nil),       // 11: google.protobuf.Empty
}
var file_swarm_grpc_service_proto_depIdxs = []int32{
	10, // 0: grpc.RegistrationRequest.Torrent:type_name -> google.protobuf.Struct
	10, // 1: grpc.RegistrationReply.Config:type_name -> google.protobuf.Struct
	10, // 2: grpc.TorrentsReply.Torrent:type_name -> google.protobuf.Struct
	11, // 3: grpc.WorkerService.GetTorrents:input_type -> google.protobuf.Empty
	3,  // 4: grpc.WorkerService.GetTorrentScore:input_type -> grpc.TorrentScoreRequest
	5,  // 5: grpc.WorkerService.DropTorrent:input_type -> grpc.TorrentDropRequest
	7,  // 6: grpc.WorkerService.SaveTorrentFile:input_type -> grpc.TFileSaveRequest
	11, // 7: grpc.WorkerService.GetSystemFreeSpace:input_type -> google.protobuf.Empty
	11, // 8: grpc.MasterService.Ping:input_type -> google.protobuf.Empty
	0,  // 9: grpc.MasterService.Register:input_type -> grpc.RegistrationRequest
	2,  // 10: grpc.WorkerService.GetTorrents:output_type -> grpc.TorrentsReply
	4,  // 11: grpc.WorkerService.GetTorrentScore:output_type -> grpc.TorrentScoreReply
	8,  // 12: grpc.WorkerService.DropTorrent:output_type -> grpc.TFileSaveReply
	8,  // 13: grpc.WorkerService.SaveTorrentFile:output_type -> grpc.TFileSaveReply
	9,  // 14: grpc.WorkerService.GetSystemFreeSpace:output_type -> grpc.SystemSpaceReply
	11, // 15: grpc.MasterService.Ping:output_type -> google.protobuf.Empty
	1,  // 16: grpc.MasterService.Register:output_type -> grpc.RegistrationReply
	10, // [10:17] is the sub-list for method output_type
	3,  // [3:10] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_swarm_grpc_service_proto_init() }
func file_swarm_grpc_service_proto_init() {
	if File_swarm_grpc_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_swarm_grpc_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistrationRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_swarm_grpc_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistrationReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_swarm_grpc_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TorrentsReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_swarm_grpc_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TorrentScoreRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_swarm_grpc_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TorrentScoreReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_swarm_grpc_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TorrentDropRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_swarm_grpc_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TorrentDropReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_swarm_grpc_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TFileSaveRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_swarm_grpc_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TFileSaveReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_swarm_grpc_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SystemSpaceReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_swarm_grpc_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_swarm_grpc_service_proto_goTypes,
		DependencyIndexes: file_swarm_grpc_service_proto_depIdxs,
		MessageInfos:      file_swarm_grpc_service_proto_msgTypes,
	}.Build()
	File_swarm_grpc_service_proto = out.File
	file_swarm_grpc_service_proto_rawDesc = nil
	file_swarm_grpc_service_proto_goTypes = nil
	file_swarm_grpc_service_proto_depIdxs = nil
}
