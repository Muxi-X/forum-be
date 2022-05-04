// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/chat.proto

package chat

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type CreateRequest struct {
	UserId               string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	TargetUserId         string   `protobuf:"bytes,2,opt,name=target_user_id,json=targetUserId,proto3" json:"target_user_id,omitempty"`
	Message              string   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed7e7dde45555b7d, []int{0}
}

func (m *CreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRequest.Unmarshal(m, b)
}
func (m *CreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRequest.Marshal(b, m, deterministic)
}
func (m *CreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRequest.Merge(m, src)
}
func (m *CreateRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRequest.Size(m)
}
func (m *CreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRequest proto.InternalMessageInfo

func (m *CreateRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *CreateRequest) GetTargetUserId() string {
	if m != nil {
		return m.TargetUserId
	}
	return ""
}

func (m *CreateRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Response struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed7e7dde45555b7d, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

type GetListRequest struct {
	UserId               string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetListRequest) Reset()         { *m = GetListRequest{} }
func (m *GetListRequest) String() string { return proto.CompactTextString(m) }
func (*GetListRequest) ProtoMessage()    {}
func (*GetListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed7e7dde45555b7d, []int{2}
}

func (m *GetListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetListRequest.Unmarshal(m, b)
}
func (m *GetListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetListRequest.Marshal(b, m, deterministic)
}
func (m *GetListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetListRequest.Merge(m, src)
}
func (m *GetListRequest) XXX_Size() int {
	return xxx_messageInfo_GetListRequest.Size(m)
}
func (m *GetListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetListRequest proto.InternalMessageInfo

func (m *GetListRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type GetListResponse struct {
	List                 []string `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetListResponse) Reset()         { *m = GetListResponse{} }
func (m *GetListResponse) String() string { return proto.CompactTextString(m) }
func (*GetListResponse) ProtoMessage()    {}
func (*GetListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ed7e7dde45555b7d, []int{3}
}

func (m *GetListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetListResponse.Unmarshal(m, b)
}
func (m *GetListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetListResponse.Marshal(b, m, deterministic)
}
func (m *GetListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetListResponse.Merge(m, src)
}
func (m *GetListResponse) XXX_Size() int {
	return xxx_messageInfo_GetListResponse.Size(m)
}
func (m *GetListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetListResponse proto.InternalMessageInfo

func (m *GetListResponse) GetList() []string {
	if m != nil {
		return m.List
	}
	return nil
}

func init() {
	proto.RegisterType((*CreateRequest)(nil), "chat.CreateRequest")
	proto.RegisterType((*Response)(nil), "chat.Response")
	proto.RegisterType((*GetListRequest)(nil), "chat.GetListRequest")
	proto.RegisterType((*GetListResponse)(nil), "chat.GetListResponse")
}

func init() { proto.RegisterFile("proto/chat.proto", fileDescriptor_ed7e7dde45555b7d) }

var fileDescriptor_ed7e7dde45555b7d = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x8d, 0x2d, 0x89, 0x1d, 0x35, 0xca, 0xa8, 0xb8, 0xf4, 0x54, 0x16, 0x85, 0x7a, 0x69,
	0x41, 0x2f, 0xde, 0x7b, 0x10, 0xc1, 0xd3, 0x8a, 0xe7, 0xb2, 0x36, 0x43, 0xb2, 0xa0, 0x26, 0xee,
	0x4c, 0xc4, 0x9f, 0x2f, 0xd9, 0x35, 0x81, 0x78, 0xf1, 0xf6, 0xe6, 0x7d, 0x33, 0xcc, 0xbc, 0x81,
	0xd3, 0xc6, 0xd7, 0x52, 0xaf, 0x77, 0x95, 0x95, 0x55, 0x90, 0x38, 0xed, 0xb4, 0xae, 0xe0, 0x78,
	0xe3, 0xc9, 0x0a, 0x19, 0xfa, 0x6c, 0x89, 0x05, 0x2f, 0x21, 0x6b, 0x99, 0xfc, 0xd6, 0x15, 0x2a,
	0x59, 0x24, 0xcb, 0x99, 0x49, 0xbb, 0xf2, 0xb1, 0xc0, 0x2b, 0xc8, 0xc5, 0xfa, 0x92, 0x64, 0xdb,
	0xf3, 0xfd, 0xc0, 0x8f, 0xa2, 0xfb, 0x12, 0xbb, 0x14, 0x64, 0xef, 0xc4, 0x6c, 0x4b, 0x52, 0x93,
	0x80, 0xfb, 0x52, 0x03, 0x1c, 0x18, 0xe2, 0xa6, 0xfe, 0x60, 0xd2, 0x37, 0x90, 0x3f, 0x90, 0x3c,
	0x39, 0x96, 0xff, 0xd6, 0xea, 0x6b, 0x38, 0x19, 0x5a, 0xe3, 0x34, 0x22, 0x4c, 0xdf, 0x1c, 0x8b,
	0x4a, 0x16, 0x93, 0xe5, 0xcc, 0x04, 0x7d, 0xfb, 0x0d, 0x87, 0x9b, 0xca, 0xca, 0x33, 0xf9, 0x2f,
	0xb7, 0x23, 0x5c, 0x43, 0x1a, 0x63, 0xe1, 0xd9, 0x2a, 0x64, 0x1e, 0x85, 0x9c, 0xe7, 0xd1, 0x1c,
	0xee, 0xd9, 0xc3, 0x7b, 0xc8, 0x7e, 0xd7, 0xe0, 0x79, 0x84, 0xe3, 0x03, 0xe7, 0x17, 0x7f, 0xdc,
	0x7e, 0xf2, 0x35, 0x0d, 0xef, 0xbc, 0xfb, 0x09, 0x00, 0x00, 0xff, 0xff, 0x3e, 0x4a, 0x00, 0xe2,
	0x62, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ChatService service

type ChatServiceClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*Response, error)
	GetList(ctx context.Context, in *GetListRequest, opts ...client.CallOption) (*GetListResponse, error)
}

type chatServiceClient struct {
	c           client.Client
	serviceName string
}

func NewChatServiceClient(serviceName string, c client.Client) ChatServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "chat"
	}
	return &chatServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *chatServiceClient) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "ChatService.Create", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetList(ctx context.Context, in *GetListRequest, opts ...client.CallOption) (*GetListResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ChatService.GetList", in)
	out := new(GetListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ChatService service

type ChatServiceHandler interface {
	Create(context.Context, *CreateRequest, *Response) error
	GetList(context.Context, *GetListRequest, *GetListResponse) error
}

func RegisterChatServiceHandler(s server.Server, hdlr ChatServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&ChatService{hdlr}, opts...))
}

type ChatService struct {
	ChatServiceHandler
}

func (h *ChatService) Create(ctx context.Context, in *CreateRequest, out *Response) error {
	return h.ChatServiceHandler.Create(ctx, in, out)
}

func (h *ChatService) GetList(ctx context.Context, in *GetListRequest, out *GetListResponse) error {
	return h.ChatServiceHandler.GetList(ctx, in, out)
}