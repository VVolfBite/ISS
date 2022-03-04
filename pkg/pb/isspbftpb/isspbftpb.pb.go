//
//Copyright IBM Corp. All Rights Reserved.
//
//SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: isspbftpb/isspbftpb.proto

package isspbftpb

import (
	requestpb "github.com/hyperledger-labs/mirbft/pkg/pb/requestpb"
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

type Preprepare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn      uint64           `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	View    uint64           `protobuf:"varint,2,opt,name=view,proto3" json:"view,omitempty"`
	Batch   *requestpb.Batch `protobuf:"bytes,3,opt,name=batch,proto3" json:"batch,omitempty"`
	Aborted bool             `protobuf:"varint,4,opt,name=aborted,proto3" json:"aborted,omitempty"`
}

func (x *Preprepare) Reset() {
	*x = Preprepare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_isspbftpb_isspbftpb_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Preprepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Preprepare) ProtoMessage() {}

func (x *Preprepare) ProtoReflect() protoreflect.Message {
	mi := &file_isspbftpb_isspbftpb_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Preprepare.ProtoReflect.Descriptor instead.
func (*Preprepare) Descriptor() ([]byte, []int) {
	return file_isspbftpb_isspbftpb_proto_rawDescGZIP(), []int{0}
}

func (x *Preprepare) GetSn() uint64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *Preprepare) GetView() uint64 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *Preprepare) GetBatch() *requestpb.Batch {
	if x != nil {
		return x.Batch
	}
	return nil
}

func (x *Preprepare) GetAborted() bool {
	if x != nil {
		return x.Aborted
	}
	return false
}

type Prepare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn     uint64 `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	View   uint64 `protobuf:"varint,2,opt,name=view,proto3" json:"view,omitempty"`
	Digest []byte `protobuf:"bytes,3,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (x *Prepare) Reset() {
	*x = Prepare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_isspbftpb_isspbftpb_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Prepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Prepare) ProtoMessage() {}

func (x *Prepare) ProtoReflect() protoreflect.Message {
	mi := &file_isspbftpb_isspbftpb_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Prepare.ProtoReflect.Descriptor instead.
func (*Prepare) Descriptor() ([]byte, []int) {
	return file_isspbftpb_isspbftpb_proto_rawDescGZIP(), []int{1}
}

func (x *Prepare) GetSn() uint64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *Prepare) GetView() uint64 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *Prepare) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

type Commit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn     uint64 `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	View   uint64 `protobuf:"varint,2,opt,name=view,proto3" json:"view,omitempty"`
	Digest []byte `protobuf:"bytes,3,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (x *Commit) Reset() {
	*x = Commit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_isspbftpb_isspbftpb_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Commit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Commit) ProtoMessage() {}

func (x *Commit) ProtoReflect() protoreflect.Message {
	mi := &file_isspbftpb_isspbftpb_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Commit.ProtoReflect.Descriptor instead.
func (*Commit) Descriptor() ([]byte, []int) {
	return file_isspbftpb_isspbftpb_proto_rawDescGZIP(), []int{2}
}

func (x *Commit) GetSn() uint64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *Commit) GetView() uint64 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *Commit) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

type ReqWaitReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sn   uint64 `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	View uint64 `protobuf:"varint,2,opt,name=view,proto3" json:"view,omitempty"`
}

func (x *ReqWaitReference) Reset() {
	*x = ReqWaitReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_isspbftpb_isspbftpb_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqWaitReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqWaitReference) ProtoMessage() {}

func (x *ReqWaitReference) ProtoReflect() protoreflect.Message {
	mi := &file_isspbftpb_isspbftpb_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqWaitReference.ProtoReflect.Descriptor instead.
func (*ReqWaitReference) Descriptor() ([]byte, []int) {
	return file_isspbftpb_isspbftpb_proto_rawDescGZIP(), []int{3}
}

func (x *ReqWaitReference) GetSn() uint64 {
	if x != nil {
		return x.Sn
	}
	return 0
}

func (x *ReqWaitReference) GetView() uint64 {
	if x != nil {
		return x.View
	}
	return 0
}

type PersistPreprepare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Preprepare *Preprepare `protobuf:"bytes,1,opt,name=preprepare,proto3" json:"preprepare,omitempty"`
}

func (x *PersistPreprepare) Reset() {
	*x = PersistPreprepare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_isspbftpb_isspbftpb_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistPreprepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistPreprepare) ProtoMessage() {}

func (x *PersistPreprepare) ProtoReflect() protoreflect.Message {
	mi := &file_isspbftpb_isspbftpb_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersistPreprepare.ProtoReflect.Descriptor instead.
func (*PersistPreprepare) Descriptor() ([]byte, []int) {
	return file_isspbftpb_isspbftpb_proto_rawDescGZIP(), []int{4}
}

func (x *PersistPreprepare) GetPreprepare() *Preprepare {
	if x != nil {
		return x.Preprepare
	}
	return nil
}

type PersistPrepare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prepare *Prepare `protobuf:"bytes,1,opt,name=prepare,proto3" json:"prepare,omitempty"`
}

func (x *PersistPrepare) Reset() {
	*x = PersistPrepare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_isspbftpb_isspbftpb_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistPrepare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistPrepare) ProtoMessage() {}

func (x *PersistPrepare) ProtoReflect() protoreflect.Message {
	mi := &file_isspbftpb_isspbftpb_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersistPrepare.ProtoReflect.Descriptor instead.
func (*PersistPrepare) Descriptor() ([]byte, []int) {
	return file_isspbftpb_isspbftpb_proto_rawDescGZIP(), []int{5}
}

func (x *PersistPrepare) GetPrepare() *Prepare {
	if x != nil {
		return x.Prepare
	}
	return nil
}

type PersistCommit struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Commit *Commit `protobuf:"bytes,1,opt,name=commit,proto3" json:"commit,omitempty"`
}

func (x *PersistCommit) Reset() {
	*x = PersistCommit{}
	if protoimpl.UnsafeEnabled {
		mi := &file_isspbftpb_isspbftpb_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistCommit) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistCommit) ProtoMessage() {}

func (x *PersistCommit) ProtoReflect() protoreflect.Message {
	mi := &file_isspbftpb_isspbftpb_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersistCommit.ProtoReflect.Descriptor instead.
func (*PersistCommit) Descriptor() ([]byte, []int) {
	return file_isspbftpb_isspbftpb_proto_rawDescGZIP(), []int{6}
}

func (x *PersistCommit) GetCommit() *Commit {
	if x != nil {
		return x.Commit
	}
	return nil
}

type PreprepareHashOrigin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Preprepare *Preprepare `protobuf:"bytes,1,opt,name=preprepare,proto3" json:"preprepare,omitempty"`
}

func (x *PreprepareHashOrigin) Reset() {
	*x = PreprepareHashOrigin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_isspbftpb_isspbftpb_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PreprepareHashOrigin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PreprepareHashOrigin) ProtoMessage() {}

func (x *PreprepareHashOrigin) ProtoReflect() protoreflect.Message {
	mi := &file_isspbftpb_isspbftpb_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PreprepareHashOrigin.ProtoReflect.Descriptor instead.
func (*PreprepareHashOrigin) Descriptor() ([]byte, []int) {
	return file_isspbftpb_isspbftpb_proto_rawDescGZIP(), []int{7}
}

func (x *PreprepareHashOrigin) GetPreprepare() *Preprepare {
	if x != nil {
		return x.Preprepare
	}
	return nil
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_isspbftpb_isspbftpb_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_isspbftpb_isspbftpb_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_isspbftpb_isspbftpb_proto_rawDescGZIP(), []int{8}
}

var File_isspbftpb_isspbftpb_proto protoreflect.FileDescriptor

var file_isspbftpb_isspbftpb_proto_rawDesc = []byte{
	0x0a, 0x19, 0x69, 0x73, 0x73, 0x70, 0x62, 0x66, 0x74, 0x70, 0x62, 0x2f, 0x69, 0x73, 0x73, 0x70,
	0x62, 0x66, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x69, 0x73, 0x73,
	0x70, 0x62, 0x66, 0x74, 0x70, 0x62, 0x1a, 0x19, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70,
	0x62, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x72, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x73, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x76,
	0x69, 0x65, 0x77, 0x12, 0x26, 0x0a, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x05, 0x62, 0x61, 0x74, 0x63, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x62, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x61, 0x62,
	0x6f, 0x72, 0x74, 0x65, 0x64, 0x22, 0x45, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x73, 0x6e,
	0x12, 0x12, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04,
	0x76, 0x69, 0x65, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x44, 0x0a, 0x06,
	0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x73, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x76, 0x69, 0x65, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x69,
	0x67, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x64, 0x69, 0x67, 0x65,
	0x73, 0x74, 0x22, 0x36, 0x0a, 0x10, 0x52, 0x65, 0x71, 0x57, 0x61, 0x69, 0x74, 0x52, 0x65, 0x66,
	0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x73, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x73, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x76, 0x69, 0x65, 0x77, 0x22, 0x4a, 0x0a, 0x11, 0x50, 0x65,
	0x72, 0x73, 0x69, 0x73, 0x74, 0x50, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12,
	0x35, 0x0a, 0x0a, 0x70, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x69, 0x73, 0x73, 0x70, 0x62, 0x66, 0x74, 0x70, 0x62, 0x2e,
	0x50, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x0a, 0x70, 0x72, 0x65, 0x70,
	0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x22, 0x3e, 0x0a, 0x0e, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73,
	0x74, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x2c, 0x0a, 0x07, 0x70, 0x72, 0x65, 0x70,
	0x61, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x69, 0x73, 0x73, 0x70,
	0x62, 0x66, 0x74, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x07, 0x70,
	0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x22, 0x3a, 0x0a, 0x0d, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73,
	0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x29, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x69, 0x73, 0x73, 0x70, 0x62, 0x66,
	0x74, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x06, 0x63, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x22, 0x4d, 0x0a, 0x14, 0x50, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x48, 0x61, 0x73, 0x68, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x35, 0x0a, 0x0a, 0x70, 0x72,
	0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x69, 0x73, 0x73, 0x70, 0x62, 0x66, 0x74, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x72,
	0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x0a, 0x70, 0x72, 0x65, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72,
	0x65, 0x22, 0x08, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x35, 0x5a, 0x33, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x79, 0x70, 0x65, 0x72, 0x6c,
	0x65, 0x64, 0x67, 0x65, 0x72, 0x2d, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x6d, 0x69, 0x72, 0x62, 0x66,
	0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x69, 0x73, 0x73, 0x70, 0x62, 0x66, 0x74,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_isspbftpb_isspbftpb_proto_rawDescOnce sync.Once
	file_isspbftpb_isspbftpb_proto_rawDescData = file_isspbftpb_isspbftpb_proto_rawDesc
)

func file_isspbftpb_isspbftpb_proto_rawDescGZIP() []byte {
	file_isspbftpb_isspbftpb_proto_rawDescOnce.Do(func() {
		file_isspbftpb_isspbftpb_proto_rawDescData = protoimpl.X.CompressGZIP(file_isspbftpb_isspbftpb_proto_rawDescData)
	})
	return file_isspbftpb_isspbftpb_proto_rawDescData
}

var file_isspbftpb_isspbftpb_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_isspbftpb_isspbftpb_proto_goTypes = []interface{}{
	(*Preprepare)(nil),           // 0: isspbftpb.Preprepare
	(*Prepare)(nil),              // 1: isspbftpb.Prepare
	(*Commit)(nil),               // 2: isspbftpb.Commit
	(*ReqWaitReference)(nil),     // 3: isspbftpb.ReqWaitReference
	(*PersistPreprepare)(nil),    // 4: isspbftpb.PersistPreprepare
	(*PersistPrepare)(nil),       // 5: isspbftpb.PersistPrepare
	(*PersistCommit)(nil),        // 6: isspbftpb.PersistCommit
	(*PreprepareHashOrigin)(nil), // 7: isspbftpb.PreprepareHashOrigin
	(*Status)(nil),               // 8: isspbftpb.Status
	(*requestpb.Batch)(nil),      // 9: requestpb.Batch
}
var file_isspbftpb_isspbftpb_proto_depIdxs = []int32{
	9, // 0: isspbftpb.Preprepare.batch:type_name -> requestpb.Batch
	0, // 1: isspbftpb.PersistPreprepare.preprepare:type_name -> isspbftpb.Preprepare
	1, // 2: isspbftpb.PersistPrepare.prepare:type_name -> isspbftpb.Prepare
	2, // 3: isspbftpb.PersistCommit.commit:type_name -> isspbftpb.Commit
	0, // 4: isspbftpb.PreprepareHashOrigin.preprepare:type_name -> isspbftpb.Preprepare
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_isspbftpb_isspbftpb_proto_init() }
func file_isspbftpb_isspbftpb_proto_init() {
	if File_isspbftpb_isspbftpb_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_isspbftpb_isspbftpb_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Preprepare); i {
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
		file_isspbftpb_isspbftpb_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Prepare); i {
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
		file_isspbftpb_isspbftpb_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Commit); i {
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
		file_isspbftpb_isspbftpb_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqWaitReference); i {
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
		file_isspbftpb_isspbftpb_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersistPreprepare); i {
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
		file_isspbftpb_isspbftpb_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersistPrepare); i {
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
		file_isspbftpb_isspbftpb_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersistCommit); i {
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
		file_isspbftpb_isspbftpb_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PreprepareHashOrigin); i {
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
		file_isspbftpb_isspbftpb_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
			RawDescriptor: file_isspbftpb_isspbftpb_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_isspbftpb_isspbftpb_proto_goTypes,
		DependencyIndexes: file_isspbftpb_isspbftpb_proto_depIdxs,
		MessageInfos:      file_isspbftpb_isspbftpb_proto_msgTypes,
	}.Build()
	File_isspbftpb_isspbftpb_proto = out.File
	file_isspbftpb_isspbftpb_proto_rawDesc = nil
	file_isspbftpb_isspbftpb_proto_goTypes = nil
	file_isspbftpb_isspbftpb_proto_depIdxs = nil
}
