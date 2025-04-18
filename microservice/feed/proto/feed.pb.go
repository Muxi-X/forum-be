// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.27.2
// source: proto/feed.proto

package proto

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

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	mi := &file_proto_feed_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feed_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_proto_feed_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type PushRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action       string  `protobuf:"bytes,1,opt,name=action,proto3" json:"action,omitempty"`
	UserId       uint32  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Source       *Source `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	TargetUserId uint32  `protobuf:"varint,4,opt,name=target_user_id,json=targetUserId,proto3" json:"target_user_id,omitempty"`
	Content      string  `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *PushRequest) Reset() {
	*x = PushRequest{}
	mi := &file_proto_feed_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PushRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushRequest) ProtoMessage() {}

func (x *PushRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feed_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushRequest.ProtoReflect.Descriptor instead.
func (*PushRequest) Descriptor() ([]byte, []int) {
	return file_proto_feed_proto_rawDescGZIP(), []int{1}
}

func (x *PushRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *PushRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *PushRequest) GetSource() *Source {
	if x != nil {
		return x.Source
	}
	return nil
}

func (x *PushRequest) GetTargetUserId() uint32 {
	if x != nil {
		return x.TargetUserId
	}
	return 0
}

func (x *PushRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Response) Reset() {
	*x = Response{}
	mi := &file_proto_feed_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feed_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_proto_feed_proto_rawDescGZIP(), []int{2}
}

type ListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LastId       uint32 `protobuf:"varint,1,opt,name=last_id,json=lastId,proto3" json:"last_id,omitempty"`
	Limit        uint32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset       uint32 `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
	UserId       uint32 `protobuf:"varint,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Pagination   bool   `protobuf:"varint,5,opt,name=pagination,proto3" json:"pagination,omitempty"`
	TargetUserId uint32 `protobuf:"varint,6,opt,name=target_user_id,json=targetUserId,proto3" json:"target_user_id,omitempty"`
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	mi := &file_proto_feed_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feed_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRequest.ProtoReflect.Descriptor instead.
func (*ListRequest) Descriptor() ([]byte, []int) {
	return file_proto_feed_proto_rawDescGZIP(), []int{3}
}

func (x *ListRequest) GetLastId() uint32 {
	if x != nil {
		return x.LastId
	}
	return 0
}

func (x *ListRequest) GetLimit() uint32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListRequest) GetOffset() uint32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *ListRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ListRequest) GetPagination() bool {
	if x != nil {
		return x.Pagination
	}
	return false
}

func (x *ListRequest) GetTargetUserId() uint32 {
	if x != nil {
		return x.TargetUserId
	}
	return 0
}

type ListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*FeedItem `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (x *ListResponse) Reset() {
	*x = ListResponse{}
	mi := &file_proto_feed_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListResponse) ProtoMessage() {}

func (x *ListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feed_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListResponse.ProtoReflect.Descriptor instead.
func (*ListResponse) Descriptor() ([]byte, []int) {
	return file_proto_feed_proto_rawDescGZIP(), []int{4}
}

func (x *ListResponse) GetList() []*FeedItem {
	if x != nil {
		return x.List
	}
	return nil
}

type FeedItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          uint32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Action      string  `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`                               // action：创建、编辑、删除、评论、加入
	ShowDivider bool    `protobuf:"varint,3,opt,name=show_divider,json=showDivider,proto3" json:"show_divider,omitempty"` // 分隔线
	CreateTime  string  `protobuf:"bytes,4,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`     // 时间
	User        *User   `protobuf:"bytes,5,opt,name=user,proto3" json:"user,omitempty"`
	Source      *Source `protobuf:"bytes,6,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *FeedItem) Reset() {
	*x = FeedItem{}
	mi := &file_proto_feed_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FeedItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FeedItem) ProtoMessage() {}

func (x *FeedItem) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feed_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FeedItem.ProtoReflect.Descriptor instead.
func (*FeedItem) Descriptor() ([]byte, []int) {
	return file_proto_feed_proto_rawDescGZIP(), []int{5}
}

func (x *FeedItem) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *FeedItem) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *FeedItem) GetShowDivider() bool {
	if x != nil {
		return x.ShowDivider
	}
	return false
}

func (x *FeedItem) GetCreateTime() string {
	if x != nil {
		return x.CreateTime
	}
	return ""
}

func (x *FeedItem) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *FeedItem) GetSource() *Source {
	if x != nil {
		return x.Source
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	AvatarUrl string `protobuf:"bytes,3,opt,name=avatar_url,json=avatarUrl,proto3" json:"avatar_url,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_proto_feed_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feed_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_proto_feed_proto_rawDescGZIP(), []int{6}
}

func (x *User) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetAvatarUrl() string {
	if x != nil {
		return x.AvatarUrl
	}
	return ""
}

// source: feed的来源.
type Source struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       uint32 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	TypeName string `protobuf:"bytes,1,opt,name=type_name,json=typeName,proto3" json:"type_name,omitempty"`
	Name     string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Source) Reset() {
	*x = Source{}
	mi := &file_proto_feed_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feed_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_proto_feed_proto_rawDescGZIP(), []int{7}
}

func (x *Source) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Source) GetTypeName() string {
	if x != nil {
		return x.TypeName
	}
	return ""
}

func (x *Source) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_proto_feed_proto protoreflect.FileDescriptor

var file_proto_feed_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x66, 0x65, 0x65, 0x64, 0x22, 0x19, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x02, 0x69, 0x64, 0x22, 0xa4, 0x01, 0x0a, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x2e, 0x53, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x0a, 0x0a, 0x08, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xb3, 0x01, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6c, 0x61, 0x73, 0x74, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x32, 0x0a, 0x0c,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x04,
	0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x65, 0x65,
	0x64, 0x2e, 0x46, 0x65, 0x65, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74,
	0x22, 0xbc, 0x01, 0x0a, 0x08, 0x46, 0x65, 0x65, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x68, 0x6f, 0x77, 0x5f, 0x64, 0x69,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x73, 0x68, 0x6f,
	0x77, 0x44, 0x69, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x06, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x66, 0x65, 0x65, 0x64,
	0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22,
	0x49, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x55, 0x72, 0x6c, 0x22, 0x49, 0x0a, 0x06, 0x53, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x79, 0x70, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0x96, 0x01, 0x0a, 0x0b, 0x46, 0x65, 0x65, 0x64, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x50, 0x75, 0x73, 0x68, 0x12, 0x11, 0x2e,
	0x66, 0x65, 0x65, 0x64, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x0e, 0x2e, 0x66, 0x65, 0x65, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x2f, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x11, 0x2e, 0x66, 0x65, 0x65,
	0x64, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x66, 0x65, 0x65, 0x64, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x0d, 0x2e,
	0x66, 0x65, 0x65, 0x64, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x66,
	0x65, 0x65, 0x64, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x09,
	0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_feed_proto_rawDescOnce sync.Once
	file_proto_feed_proto_rawDescData = file_proto_feed_proto_rawDesc
)

func file_proto_feed_proto_rawDescGZIP() []byte {
	file_proto_feed_proto_rawDescOnce.Do(func() {
		file_proto_feed_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_feed_proto_rawDescData)
	})
	return file_proto_feed_proto_rawDescData
}

var file_proto_feed_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_feed_proto_goTypes = []any{
	(*Request)(nil),      // 0: feed.Request
	(*PushRequest)(nil),  // 1: feed.PushRequest
	(*Response)(nil),     // 2: feed.Response
	(*ListRequest)(nil),  // 3: feed.ListRequest
	(*ListResponse)(nil), // 4: feed.ListResponse
	(*FeedItem)(nil),     // 5: feed.FeedItem
	(*User)(nil),         // 6: feed.User
	(*Source)(nil),       // 7: feed.Source
}
var file_proto_feed_proto_depIdxs = []int32{
	7, // 0: feed.PushRequest.source:type_name -> feed.Source
	5, // 1: feed.ListResponse.list:type_name -> feed.FeedItem
	6, // 2: feed.FeedItem.user:type_name -> feed.User
	7, // 3: feed.FeedItem.source:type_name -> feed.Source
	1, // 4: feed.FeedService.Push:input_type -> feed.PushRequest
	3, // 5: feed.FeedService.List:input_type -> feed.ListRequest
	0, // 6: feed.FeedService.Delete:input_type -> feed.Request
	2, // 7: feed.FeedService.Push:output_type -> feed.Response
	4, // 8: feed.FeedService.List:output_type -> feed.ListResponse
	2, // 9: feed.FeedService.Delete:output_type -> feed.Response
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_feed_proto_init() }
func file_proto_feed_proto_init() {
	if File_proto_feed_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_feed_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_feed_proto_goTypes,
		DependencyIndexes: file_proto_feed_proto_depIdxs,
		MessageInfos:      file_proto_feed_proto_msgTypes,
	}.Build()
	File_proto_feed_proto = out.File
	file_proto_feed_proto_rawDesc = nil
	file_proto_feed_proto_goTypes = nil
	file_proto_feed_proto_depIdxs = nil
}
