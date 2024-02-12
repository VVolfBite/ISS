// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0-devel
// 	protoc        v3.12.4
// source: pbftorderer.proto

package protobufs

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PbftPreprepare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn      int32  `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	View    int32  `protobuf:"varint,2,opt,name=view,proto3" json:"view,omitempty"`
	Leader  int32  `protobuf:"varint,3,opt,name=leader,proto3" json:"leader,omitempty"`
	Batch   *Batch `protobuf:"bytes,4,opt,name=batch,proto3" json:"batch,omitempty"`
	Aborted bool   `protobuf:"varint,5,opt,name=aborted,proto3" json:"aborted,omitempty"`
	Ts      int64  `protobuf:"varint,6,opt,name=ts,proto3" json:"ts,omitempty"` // Timestamp to be set by the receiver at message reception.
}

func (x *PbftPreprepare) Reset() {
	*x = PbftPreprepare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbftorderer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbftPreprepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbftPreprepare) ProtoMessage() {}

func (x *PbftPreprepare) ProtoReflect() protoreflect.Message {
	mi := &file_pbftorderer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbftPreprepare.ProtoReflect.Descriptor instead.
func (*PbftPreprepare) Descriptor() ([]byte, []int) {
	return file_pbftorderer_proto_rawDescGZIP(), []int{0}
}

func (x *PbftPreprepare) GetSn() int32 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *PbftPreprepare) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *PbftPreprepare) GetLeader() int32 {
	if x != nil {
		return x.Leader
	}
	return 0
}

func (x *PbftPreprepare) GetBatch() *Batch {
	if x != nil {
		return x.Batch
	}
	return nil
}

func (x *PbftPreprepare) GetAborted() bool {
	if x != nil {
		return x.Aborted
	}
	return false
}

func (x *PbftPreprepare) GetTs() int64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

type PbftPrepare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn     int32  `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	View   int32  `protobuf:"varint,2,opt,name=view,proto3" json:"view,omitempty"`
	Digest []byte `protobuf:"bytes,3,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (x *PbftPrepare) Reset() {
	*x = PbftPrepare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbftorderer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbftPrepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbftPrepare) ProtoMessage() {}

func (x *PbftPrepare) ProtoReflect() protoreflect.Message {
	mi := &file_pbftorderer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbftPrepare.ProtoReflect.Descriptor instead.
func (*PbftPrepare) Descriptor() ([]byte, []int) {
	return file_pbftorderer_proto_rawDescGZIP(), []int{1}
}

func (x *PbftPrepare) GetSn() int32 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *PbftPrepare) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *PbftPrepare) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

type PbftCommit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn     int32  `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	View   int32  `protobuf:"varint,2,opt,name=view,proto3" json:"view,omitempty"`
	Digest []byte `protobuf:"bytes,3,opt,name=digest,proto3" json:"digest,omitempty"`
	Ts     int64  `protobuf:"varint,4,opt,name=ts,proto3" json:"ts,omitempty"` // Timestamp to be set by the receiver at message reception.
}

func (x *PbftCommit) Reset() {
	*x = PbftCommit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbftorderer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbftCommit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbftCommit) ProtoMessage() {}

func (x *PbftCommit) ProtoReflect() protoreflect.Message {
	mi := &file_pbftorderer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbftCommit.ProtoReflect.Descriptor instead.
func (*PbftCommit) Descriptor() ([]byte, []int) {
	return file_pbftorderer_proto_rawDescGZIP(), []int{2}
}

func (x *PbftCommit) GetSn() int32 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *PbftCommit) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *PbftCommit) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

func (x *PbftCommit) GetTs() int64 {
	if x != nil {
		return x.Ts
	}
	return 0
}

type PbftCheckpoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Digests [][]byte `protobuf:"bytes,1,rep,name=digests,proto3" json:"digests,omitempty"`
}

func (x *PbftCheckpoint) Reset() {
	*x = PbftCheckpoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbftorderer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbftCheckpoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbftCheckpoint) ProtoMessage() {}

func (x *PbftCheckpoint) ProtoReflect() protoreflect.Message {
	mi := &file_pbftorderer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbftCheckpoint.ProtoReflect.Descriptor instead.
func (*PbftCheckpoint) Descriptor() ([]byte, []int) {
	return file_pbftorderer_proto_rawDescGZIP(), []int{3}
}

func (x *PbftCheckpoint) GetDigests() [][]byte {
	if x != nil {
		return x.Digests
	}
	return nil
}

// Not sent over the networ, only used internally by the PBFT instance
type PbftCatchUp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PbftCatchUp) Reset() {
	*x = PbftCatchUp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbftorderer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbftCatchUp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbftCatchUp) ProtoMessage() {}

func (x *PbftCatchUp) ProtoReflect() protoreflect.Message {
	mi := &file_pbftorderer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbftCatchUp.ProtoReflect.Descriptor instead.
func (*PbftCatchUp) Descriptor() ([]byte, []int) {
	return file_pbftorderer_proto_rawDescGZIP(), []int{4}
}

type PbftViewChange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	View     int32                  `protobuf:"varint,1,opt,name=view,proto3" json:"view,omitempty"`                                                                                         // new view
	H        int32                  `protobuf:"varint,2,opt,name=h,proto3" json:"h,omitempty"`                                                                                               // latest stable checkpoitn
	Pset     map[int32]*PbftPrepare `protobuf:"bytes,3,rep,name=pset,proto3" json:"pset,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // prepared requests at previous views
	Qset     map[int32]*PbftPrepare `protobuf:"bytes,4,rep,name=qset,proto3" json:"qset,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"` // preprepared requests at previous views
	Cset     []*CheckpointMsg       `protobuf:"bytes,5,rep,name=cset,proto3" json:"cset,omitempty"`                                                                                          // all available checkpoints
	SenderId int32                  `protobuf:"varint,6,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`                                                                 // Sender ID for convenience (not strictly needed, since it's part of the ProtocolMessage already)
}

func (x *PbftViewChange) Reset() {
	*x = PbftViewChange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbftorderer_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbftViewChange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbftViewChange) ProtoMessage() {}

func (x *PbftViewChange) ProtoReflect() protoreflect.Message {
	mi := &file_pbftorderer_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbftViewChange.ProtoReflect.Descriptor instead.
func (*PbftViewChange) Descriptor() ([]byte, []int) {
	return file_pbftorderer_proto_rawDescGZIP(), []int{5}
}

func (x *PbftViewChange) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *PbftViewChange) GetH() int32 {
	if x != nil {
		return x.H
	}
	return 0
}

func (x *PbftViewChange) GetPset() map[int32]*PbftPrepare {
	if x != nil {
		return x.Pset
	}
	return nil
}

func (x *PbftViewChange) GetQset() map[int32]*PbftPrepare {
	if x != nil {
		return x.Qset
	}
	return nil
}

func (x *PbftViewChange) GetCset() []*CheckpointMsg {
	if x != nil {
		return x.Cset
	}
	return nil
}

func (x *PbftViewChange) GetSenderId() int32 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

type PbftMissingPreprepareRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	View int32 `protobuf:"varint,1,opt,name=view,proto3" json:"view,omitempty"`
}

func (x *PbftMissingPreprepareRequest) Reset() {
	*x = PbftMissingPreprepareRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbftorderer_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbftMissingPreprepareRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbftMissingPreprepareRequest) ProtoMessage() {}

func (x *PbftMissingPreprepareRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pbftorderer_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbftMissingPreprepareRequest.ProtoReflect.Descriptor instead.
func (*PbftMissingPreprepareRequest) Descriptor() ([]byte, []int) {
	return file_pbftorderer_proto_rawDescGZIP(), []int{6}
}

func (x *PbftMissingPreprepareRequest) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

type PbftMissingPreprepare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Preprepare *PbftPreprepare `protobuf:"bytes,1,opt,name=preprepare,proto3" json:"preprepare,omitempty"`
}

func (x *PbftMissingPreprepare) Reset() {
	*x = PbftMissingPreprepare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbftorderer_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbftMissingPreprepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbftMissingPreprepare) ProtoMessage() {}

func (x *PbftMissingPreprepare) ProtoReflect() protoreflect.Message {
	mi := &file_pbftorderer_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbftMissingPreprepare.ProtoReflect.Descriptor instead.
func (*PbftMissingPreprepare) Descriptor() ([]byte, []int) {
	return file_pbftorderer_proto_rawDescGZIP(), []int{7}
}

func (x *PbftMissingPreprepare) GetPreprepare() *PbftPreprepare {
	if x != nil {
		return x.Preprepare
	}
	return nil
}

type PbftNewView struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	View       int32                     `protobuf:"varint,1,opt,name=view,proto3" json:"view,omitempty"`
	Vset       map[int32]*SignedMsg      `protobuf:"bytes,2,rep,name=vset,proto3" json:"vset,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Xset       map[int32]*PbftPreprepare `protobuf:"bytes,3,rep,name=xset,proto3" json:"xset,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Checkpoint *CheckpointMsg            `protobuf:"bytes,4,opt,name=checkpoint,proto3" json:"checkpoint,omitempty"`
}

func (x *PbftNewView) Reset() {
	*x = PbftNewView{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pbftorderer_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PbftNewView) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PbftNewView) ProtoMessage() {}

func (x *PbftNewView) ProtoReflect() protoreflect.Message {
	mi := &file_pbftorderer_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PbftNewView.ProtoReflect.Descriptor instead.
func (*PbftNewView) Descriptor() ([]byte, []int) {
	return file_pbftorderer_proto_rawDescGZIP(), []int{8}
}

func (x *PbftNewView) GetView() int32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *PbftNewView) GetVset() map[int32]*SignedMsg {
	if x != nil {
		return x.Vset
	}
	return nil
}

func (x *PbftNewView) GetXset() map[int32]*PbftPreprepare {
	if x != nil {
		return x.Xset
	}
	return nil
}

func (x *PbftNewView) GetCheckpoint() *CheckpointMsg {
	if x != nil {
		return x.Checkpoint
	}
	return nil
}

var File_pbftorderer_proto protoreflect.FileDescriptor

var file_pbftorderer_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x62, 0x66, 0x74, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x1a, 0x0d,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x63,
	0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9e, 0x01,
	0x0a, 0x0e, 0x50, 0x62, 0x66, 0x74, 0x50, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x73, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x76, 0x69, 0x65, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x26, 0x0a, 0x05,
	0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x05, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x62, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x61, 0x62, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x12, 0x0e,
	0x0a, 0x02, 0x74, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x74, 0x73, 0x22, 0x49,
	0x0a, 0x0b, 0x50, 0x62, 0x66, 0x74, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x73, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x76, 0x69, 0x65,
	0x77, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x58, 0x0a, 0x0a, 0x50, 0x62, 0x66,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x02, 0x73, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x76, 0x69, 0x65, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x64,
	0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x64, 0x69, 0x67,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x74, 0x73, 0x22, 0x2a, 0x0a, 0x0e, 0x50, 0x62, 0x66, 0x74, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x07, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x73, 0x22,
	0x0d, 0x0a, 0x0b, 0x50, 0x62, 0x66, 0x74, 0x43, 0x61, 0x74, 0x63, 0x68, 0x55, 0x70, 0x22, 0x91,
	0x03, 0x0a, 0x0e, 0x50, 0x62, 0x66, 0x74, 0x56, 0x69, 0x65, 0x77, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x76, 0x69, 0x65, 0x77, 0x12, 0x0c, 0x0a, 0x01, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x01, 0x68, 0x12, 0x37, 0x0a, 0x04, 0x70, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x50, 0x62,
	0x66, 0x74, 0x56, 0x69, 0x65, 0x77, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x2e, 0x50, 0x73, 0x65,
	0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x70, 0x73, 0x65, 0x74, 0x12, 0x37, 0x0a, 0x04,
	0x71, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x50, 0x62, 0x66, 0x74, 0x56, 0x69, 0x65, 0x77, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x2e, 0x51, 0x73, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x04, 0x71, 0x73, 0x65, 0x74, 0x12, 0x2c, 0x0a, 0x04, 0x63, 0x73, 0x65, 0x74, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e,
	0x43, 0x68, 0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x4d, 0x73, 0x67, 0x52, 0x04, 0x63,
	0x73, 0x65, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64,
	0x1a, 0x4f, 0x0a, 0x09, 0x50, 0x73, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x2c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x50, 0x62, 0x66, 0x74, 0x50,
	0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x1a, 0x4f, 0x0a, 0x09, 0x51, 0x73, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x50, 0x62, 0x66, 0x74,
	0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x32, 0x0a, 0x1c, 0x50, 0x62, 0x66, 0x74, 0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e,
	0x67, 0x50, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x76, 0x69, 0x65, 0x77, 0x22, 0x52, 0x0a, 0x15, 0x50, 0x62, 0x66, 0x74, 0x4d, 0x69,
	0x73, 0x73, 0x69, 0x6e, 0x67, 0x50, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12,
	0x39, 0x0a, 0x0a, 0x70, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e,
	0x50, 0x62, 0x66, 0x74, 0x50, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x0a,
	0x70, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x22, 0xea, 0x02, 0x0a, 0x0b, 0x50,
	0x62, 0x66, 0x74, 0x4e, 0x65, 0x77, 0x56, 0x69, 0x65, 0x77, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x69,
	0x65, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x76, 0x69, 0x65, 0x77, 0x12, 0x34,
	0x0a, 0x04, 0x76, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x50, 0x62, 0x66, 0x74, 0x4e, 0x65, 0x77,
	0x56, 0x69, 0x65, 0x77, 0x2e, 0x56, 0x73, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04,
	0x76, 0x73, 0x65, 0x74, 0x12, 0x34, 0x0a, 0x04, 0x78, 0x73, 0x65, 0x74, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x50,
	0x62, 0x66, 0x74, 0x4e, 0x65, 0x77, 0x56, 0x69, 0x65, 0x77, 0x2e, 0x58, 0x73, 0x65, 0x74, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x78, 0x73, 0x65, 0x74, 0x12, 0x38, 0x0a, 0x0a, 0x63, 0x68,
	0x65, 0x63, 0x6b, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x4d, 0x73, 0x67, 0x52, 0x0a, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x1a, 0x4d, 0x0a, 0x09, 0x56, 0x73, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x53,
	0x69, 0x67, 0x6e, 0x65, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x1a, 0x52, 0x0a, 0x09, 0x58, 0x73, 0x65, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x2f, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x2e, 0x50, 0x62,
	0x66, 0x74, 0x50, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0e, 0x5a, 0x0c, 0x2e, 0x2f, 0x3b, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pbftorderer_proto_rawDescOnce sync.Once
	file_pbftorderer_proto_rawDescData = file_pbftorderer_proto_rawDesc
)

func file_pbftorderer_proto_rawDescGZIP() []byte {
	file_pbftorderer_proto_rawDescOnce.Do(func() {
		file_pbftorderer_proto_rawDescData = protoimpl.X.CompressGZIP(file_pbftorderer_proto_rawDescData)
	})
	return file_pbftorderer_proto_rawDescData
}

var file_pbftorderer_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_pbftorderer_proto_goTypes = []interface{}{
	(*PbftPreprepare)(nil),               // 0: protobufs.PbftPreprepare
	(*PbftPrepare)(nil),                  // 1: protobufs.PbftPrepare
	(*PbftCommit)(nil),                   // 2: protobufs.PbftCommit
	(*PbftCheckpoint)(nil),               // 3: protobufs.PbftCheckpoint
	(*PbftCatchUp)(nil),                  // 4: protobufs.PbftCatchUp
	(*PbftViewChange)(nil),               // 5: protobufs.PbftViewChange
	(*PbftMissingPreprepareRequest)(nil), // 6: protobufs.PbftMissingPreprepareRequest
	(*PbftMissingPreprepare)(nil),        // 7: protobufs.PbftMissingPreprepare
	(*PbftNewView)(nil),                  // 8: protobufs.PbftNewView
	nil,                                  // 9: protobufs.PbftViewChange.PsetEntry
	nil,                                  // 10: protobufs.PbftViewChange.QsetEntry
	nil,                                  // 11: protobufs.PbftNewView.VsetEntry
	nil,                                  // 12: protobufs.PbftNewView.XsetEntry
	(*Batch)(nil),                        // 13: protobufs.Batch
	(*CheckpointMsg)(nil),                // 14: protobufs.CheckpointMsg
	(*SignedMsg)(nil),                    // 15: protobufs.SignedMsg
}
var file_pbftorderer_proto_depIdxs = []int32{
	13, // 0: protobufs.PbftPreprepare.batch:type_name -> protobufs.Batch
	9,  // 1: protobufs.PbftViewChange.pset:type_name -> protobufs.PbftViewChange.PsetEntry
	10, // 2: protobufs.PbftViewChange.qset:type_name -> protobufs.PbftViewChange.QsetEntry
	14, // 3: protobufs.PbftViewChange.cset:type_name -> protobufs.CheckpointMsg
	0,  // 4: protobufs.PbftMissingPreprepare.preprepare:type_name -> protobufs.PbftPreprepare
	11, // 5: protobufs.PbftNewView.vset:type_name -> protobufs.PbftNewView.VsetEntry
	12, // 6: protobufs.PbftNewView.xset:type_name -> protobufs.PbftNewView.XsetEntry
	14, // 7: protobufs.PbftNewView.checkpoint:type_name -> protobufs.CheckpointMsg
	1,  // 8: protobufs.PbftViewChange.PsetEntry.value:type_name -> protobufs.PbftPrepare
	1,  // 9: protobufs.PbftViewChange.QsetEntry.value:type_name -> protobufs.PbftPrepare
	15, // 10: protobufs.PbftNewView.VsetEntry.value:type_name -> protobufs.SignedMsg
	0,  // 11: protobufs.PbftNewView.XsetEntry.value:type_name -> protobufs.PbftPreprepare
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_pbftorderer_proto_init() }
func file_pbftorderer_proto_init() {
	if File_pbftorderer_proto != nil {
		return
	}
	file_request_proto_init()
	file_checkpoint_proto_init()
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pbftorderer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbftPreprepare); i {
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
		file_pbftorderer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbftPrepare); i {
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
		file_pbftorderer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbftCommit); i {
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
		file_pbftorderer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbftCheckpoint); i {
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
		file_pbftorderer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbftCatchUp); i {
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
		file_pbftorderer_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbftViewChange); i {
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
		file_pbftorderer_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbftMissingPreprepareRequest); i {
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
		file_pbftorderer_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbftMissingPreprepare); i {
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
		file_pbftorderer_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PbftNewView); i {
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
			RawDescriptor: file_pbftorderer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pbftorderer_proto_goTypes,
		DependencyIndexes: file_pbftorderer_proto_depIdxs,
		MessageInfos:      file_pbftorderer_proto_msgTypes,
	}.Build()
	File_pbftorderer_proto = out.File
	file_pbftorderer_proto_rawDesc = nil
	file_pbftorderer_proto_goTypes = nil
	file_pbftorderer_proto_depIdxs = nil
}
