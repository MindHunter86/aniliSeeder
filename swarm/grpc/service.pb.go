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

type InitReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// auth phase
	WorkerId string `protobuf:"bytes,1,opt,name=WorkerId,proto3" json:"WorkerId,omitempty"`
	// registration phase
	WorkerVersion string             `protobuf:"bytes,2,opt,name=WorkerVersion,proto3" json:"WorkerVersion,omitempty"`
	WDFreeSpace   uint64             `protobuf:"varint,3,opt,name=WDFreeSpace,proto3" json:"WDFreeSpace,omitempty"`
	Torrent       []*structpb.Struct `protobuf:"bytes,4,rep,name=Torrent,proto3" json:"Torrent,omitempty"`
}

func (x *InitReply) Reset() {
	*x = InitReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitReply) ProtoMessage() {}

func (x *InitReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use InitReply.ProtoReflect.Descriptor instead.
func (*InitReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{0}
}

func (x *InitReply) GetWorkerId() string {
	if x != nil {
		return x.WorkerId
	}
	return ""
}

func (x *InitReply) GetWorkerVersion() string {
	if x != nil {
		return x.WorkerVersion
	}
	return ""
}

func (x *InitReply) GetWDFreeSpace() uint64 {
	if x != nil {
		return x.WDFreeSpace
	}
	return 0
}

func (x *InitReply) GetTorrent() []*structpb.Struct {
	if x != nil {
		return x.Torrent
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
		mi := &file_swarm_grpc_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentsReply) ProtoMessage() {}

func (x *TorrentsReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TorrentsReply.ProtoReflect.Descriptor instead.
func (*TorrentsReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{1}
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
		mi := &file_swarm_grpc_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentScoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentScoreRequest) ProtoMessage() {}

func (x *TorrentScoreRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TorrentScoreRequest.ProtoReflect.Descriptor instead.
func (*TorrentScoreRequest) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{2}
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
	Score float64 `protobuf:"fixed64,2,opt,name=Score,proto3" json:"Score,omitempty"`
}

func (x *TorrentScoreReply) Reset() {
	*x = TorrentScoreReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentScoreReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentScoreReply) ProtoMessage() {}

func (x *TorrentScoreReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TorrentScoreReply.ProtoReflect.Descriptor instead.
func (*TorrentScoreReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{3}
}

func (x *TorrentScoreReply) GetRatio() float32 {
	if x != nil {
		return x.Ratio
	}
	return 0
}

func (x *TorrentScoreReply) GetScore() float64 {
	if x != nil {
		return x.Score
	}
	return 0
}

type TorrentDropRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash     string `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Name     string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	WithData bool   `protobuf:"varint,3,opt,name=WithData,proto3" json:"WithData,omitempty"`
}

func (x *TorrentDropRequest) Reset() {
	*x = TorrentDropRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentDropRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentDropRequest) ProtoMessage() {}

func (x *TorrentDropRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TorrentDropRequest.ProtoReflect.Descriptor instead.
func (*TorrentDropRequest) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{4}
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

func (x *TorrentDropRequest) GetWithData() bool {
	if x != nil {
		return x.WithData
	}
	return false
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
		mi := &file_swarm_grpc_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentDropReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentDropReply) ProtoMessage() {}

func (x *TorrentDropReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TorrentDropReply.ProtoReflect.Descriptor instead.
func (*TorrentDropReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{5}
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

type TorrentUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Hash    string `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
	NewHash string `protobuf:"bytes,3,opt,name=NewHash,proto3" json:"NewHash,omitempty"`
	TFile   []byte `protobuf:"bytes,4,opt,name=TFile,proto3" json:"TFile,omitempty"`
}

func (x *TorrentUpdateRequest) Reset() {
	*x = TorrentUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentUpdateRequest) ProtoMessage() {}

func (x *TorrentUpdateRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TorrentUpdateRequest.ProtoReflect.Descriptor instead.
func (*TorrentUpdateRequest) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{6}
}

func (x *TorrentUpdateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TorrentUpdateRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *TorrentUpdateRequest) GetNewHash() string {
	if x != nil {
		return x.NewHash
	}
	return ""
}

func (x *TorrentUpdateRequest) GetTFile() []byte {
	if x != nil {
		return x.TFile
	}
	return nil
}

type TorrentUpdateReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	NewHash  string   `protobuf:"bytes,2,opt,name=NewHash,proto3" json:"NewHash,omitempty"`
	NewFiles []string `protobuf:"bytes,3,rep,name=NewFiles,proto3" json:"NewFiles,omitempty"`
}

func (x *TorrentUpdateReply) Reset() {
	*x = TorrentUpdateReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TorrentUpdateReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TorrentUpdateReply) ProtoMessage() {}

func (x *TorrentUpdateReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TorrentUpdateReply.ProtoReflect.Descriptor instead.
func (*TorrentUpdateReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{7}
}

func (x *TorrentUpdateReply) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TorrentUpdateReply) GetNewHash() string {
	if x != nil {
		return x.NewHash
	}
	return ""
}

func (x *TorrentUpdateReply) GetNewFiles() []string {
	if x != nil {
		return x.NewFiles
	}
	return nil
}

type TFileSaveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filename string `protobuf:"bytes,1,opt,name=Filename,proto3" json:"Filename,omitempty"`
	Payload  []byte `protobuf:"bytes,2,opt,name=Payload,proto3" json:"Payload,omitempty"`
}

func (x *TFileSaveRequest) Reset() {
	*x = TFileSaveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TFileSaveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TFileSaveRequest) ProtoMessage() {}

func (x *TFileSaveRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TFileSaveRequest.ProtoReflect.Descriptor instead.
func (*TFileSaveRequest) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{8}
}

func (x *TFileSaveRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
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

	WrittenBytes int64 `protobuf:"varint,1,opt,name=WrittenBytes,proto3" json:"WrittenBytes,omitempty"`
}

func (x *TFileSaveReply) Reset() {
	*x = TFileSaveReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_swarm_grpc_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TFileSaveReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TFileSaveReply) ProtoMessage() {}

func (x *TFileSaveReply) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TFileSaveReply.ProtoReflect.Descriptor instead.
func (*TFileSaveReply) Descriptor() ([]byte, []int) {
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{9}
}

func (x *TFileSaveReply) GetWrittenBytes() int64 {
	if x != nil {
		return x.WrittenBytes
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
		mi := &file_swarm_grpc_service_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SystemSpaceReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SystemSpaceReply) ProtoMessage() {}

func (x *SystemSpaceReply) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[10]
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
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{10}
}

func (x *SystemSpaceReply) GetFreeSpace() uint64 {
	if x != nil {
		return x.FreeSpace
	}
	return 0
}

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
		mi := &file_swarm_grpc_service_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationRequest) ProtoMessage() {}

func (x *RegistrationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_swarm_grpc_service_proto_msgTypes[11]
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
	return file_swarm_grpc_service_proto_rawDescGZIP(), []int{11}
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

var File_swarm_grpc_service_proto protoreflect.FileDescriptor

var file_swarm_grpc_service_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x77, 0x61, 0x72, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x72, 0x70, 0x63,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa2, 0x01, 0x0a, 0x09,
	0x49, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x57, 0x6f, 0x72,
	0x6b, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x57, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x57,
	0x44, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0b, 0x57, 0x44, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x12, 0x31, 0x0a,
	0x07, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x22, 0x42, 0x0a, 0x0d, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x31, 0x0a, 0x07, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x54, 0x6f, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x22, 0x3d, 0x0a, 0x13, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53,
	0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x48,
	0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x3f, 0x0a, 0x11, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x63,
	0x6f, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x61, 0x74, 0x69,
	0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x52, 0x61, 0x74, 0x69, 0x6f, 0x12, 0x14,
	0x0a, 0x05, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x53,
	0x63, 0x6f, 0x72, 0x65, 0x22, 0x58, 0x0a, 0x12, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x44,
	0x72, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61,
	0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12,
	0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x57, 0x69, 0x74, 0x68, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x57, 0x69, 0x74, 0x68, 0x44, 0x61, 0x74, 0x61, 0x22, 0x50,
	0x0a, 0x10, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x46, 0x72, 0x65, 0x65, 0x64, 0x53, 0x70, 0x61, 0x63, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x46, 0x72, 0x65, 0x65, 0x64, 0x53, 0x70, 0x61,
	0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65,
	0x22, 0x6e, 0x0a, 0x14, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68,
	0x12, 0x18, 0x0a, 0x07, 0x4e, 0x65, 0x77, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x4e, 0x65, 0x77, 0x48, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x46,
	0x69, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x54, 0x46, 0x69, 0x6c, 0x65,
	0x22, 0x5e, 0x0a, 0x12, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4e, 0x65,
	0x77, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4e, 0x65, 0x77,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x1a, 0x0a, 0x08, 0x4e, 0x65, 0x77, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x4e, 0x65, 0x77, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x22, 0x48, 0x0a, 0x10, 0x54, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x34, 0x0a, 0x0e, 0x54, 0x46,
	0x69, 0x6c, 0x65, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x22, 0x0a, 0x0c,
	0x57, 0x72, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0c, 0x57, 0x72, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x42, 0x79, 0x74, 0x65, 0x73,
	0x22, 0x30, 0x0a, 0x10, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53, 0x70, 0x61, 0x63, 0x65, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61,
	0x63, 0x65, 0x22, 0x90, 0x01, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x57, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x0a, 0x0b, 0x57, 0x44, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x57, 0x44, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61,
	0x63, 0x65, 0x12, 0x31, 0x0a, 0x07, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x54, 0x6f,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x32, 0xdf, 0x04, 0x0a, 0x0d, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x31, 0x0a, 0x04, 0x49, 0x6e, 0x69, 0x74, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x49,
	0x6e, 0x69, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x04, 0x50, 0x69,
	0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x22, 0x00, 0x12, 0x3c, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x13, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x00, 0x12, 0x47, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x19, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x6f, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x17, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x53,
	0x63, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0b, 0x44,
	0x72, 0x6f, 0x70, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x6f, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x44, 0x72, 0x6f, 0x70, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x47,
	0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x12,
	0x1a, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x54, 0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0f, 0x53, 0x61, 0x76, 0x65, 0x54,
	0x6f, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x54, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x14, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x54, 0x46, 0x69, 0x6c, 0x65, 0x53,
	0x61, 0x76, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x12, 0x47, 0x65,
	0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x46, 0x72, 0x65, 0x65, 0x53, 0x70, 0x61, 0x63, 0x65,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53, 0x70, 0x61, 0x63, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x00, 0x12, 0x43, 0x0a, 0x0f, 0x46, 0x6f, 0x72, 0x63, 0x65, 0x52, 0x65, 0x61, 0x6e, 0x6e,
	0x6f, 0x75, 0x6e, 0x63, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75,
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

var file_swarm_grpc_service_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_swarm_grpc_service_proto_goTypes = []interface{}{
	(*InitReply)(nil),            // 0: grpc.InitReply
	(*TorrentsReply)(nil),        // 1: grpc.TorrentsReply
	(*TorrentScoreRequest)(nil),  // 2: grpc.TorrentScoreRequest
	(*TorrentScoreReply)(nil),    // 3: grpc.TorrentScoreReply
	(*TorrentDropRequest)(nil),   // 4: grpc.TorrentDropRequest
	(*TorrentDropReply)(nil),     // 5: grpc.TorrentDropReply
	(*TorrentUpdateRequest)(nil), // 6: grpc.TorrentUpdateRequest
	(*TorrentUpdateReply)(nil),   // 7: grpc.TorrentUpdateReply
	(*TFileSaveRequest)(nil),     // 8: grpc.TFileSaveRequest
	(*TFileSaveReply)(nil),       // 9: grpc.TFileSaveReply
	(*SystemSpaceReply)(nil),     // 10: grpc.SystemSpaceReply
	(*RegistrationRequest)(nil),  // 11: grpc.RegistrationRequest
	(*structpb.Struct)(nil),      // 12: google.protobuf.Struct
	(*emptypb.Empty)(nil),        // 13: google.protobuf.Empty
}
var file_swarm_grpc_service_proto_depIdxs = []int32{
	12, // 0: grpc.InitReply.Torrent:type_name -> google.protobuf.Struct
	12, // 1: grpc.TorrentsReply.Torrent:type_name -> google.protobuf.Struct
	12, // 2: grpc.RegistrationRequest.Torrent:type_name -> google.protobuf.Struct
	13, // 3: grpc.WorkerService.Init:input_type -> google.protobuf.Empty
	13, // 4: grpc.WorkerService.Ping:input_type -> google.protobuf.Empty
	13, // 5: grpc.WorkerService.GetTorrents:input_type -> google.protobuf.Empty
	2,  // 6: grpc.WorkerService.GetTorrentScore:input_type -> grpc.TorrentScoreRequest
	4,  // 7: grpc.WorkerService.DropTorrent:input_type -> grpc.TorrentDropRequest
	6,  // 8: grpc.WorkerService.UpdateTorrent:input_type -> grpc.TorrentUpdateRequest
	8,  // 9: grpc.WorkerService.SaveTorrentFile:input_type -> grpc.TFileSaveRequest
	13, // 10: grpc.WorkerService.GetSystemFreeSpace:input_type -> google.protobuf.Empty
	13, // 11: grpc.WorkerService.ForceReannounce:input_type -> google.protobuf.Empty
	0,  // 12: grpc.WorkerService.Init:output_type -> grpc.InitReply
	13, // 13: grpc.WorkerService.Ping:output_type -> google.protobuf.Empty
	1,  // 14: grpc.WorkerService.GetTorrents:output_type -> grpc.TorrentsReply
	3,  // 15: grpc.WorkerService.GetTorrentScore:output_type -> grpc.TorrentScoreReply
	5,  // 16: grpc.WorkerService.DropTorrent:output_type -> grpc.TorrentDropReply
	7,  // 17: grpc.WorkerService.UpdateTorrent:output_type -> grpc.TorrentUpdateReply
	9,  // 18: grpc.WorkerService.SaveTorrentFile:output_type -> grpc.TFileSaveReply
	10, // 19: grpc.WorkerService.GetSystemFreeSpace:output_type -> grpc.SystemSpaceReply
	13, // 20: grpc.WorkerService.ForceReannounce:output_type -> google.protobuf.Empty
	12, // [12:21] is the sub-list for method output_type
	3,  // [3:12] is the sub-list for method input_type
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
			switch v := v.(*InitReply); i {
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
		file_swarm_grpc_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_swarm_grpc_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_swarm_grpc_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_swarm_grpc_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_swarm_grpc_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TorrentUpdateRequest); i {
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
			switch v := v.(*TorrentUpdateReply); i {
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
		file_swarm_grpc_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
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
		file_swarm_grpc_service_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
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
		file_swarm_grpc_service_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_swarm_grpc_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
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
