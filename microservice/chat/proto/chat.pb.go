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
	UserId               uint32   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
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

func (m *CreateRequest) GetUserId() uint32 {
	if m != nil {
		return m.UserId
	}
	return 0
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
	UserId               uint32   `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
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

func (m *GetListRequest) GetUserId() uint32 {
	if m != nil {
		return m.UserId
	}
	return 0
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
	// 230 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x50, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x75, 0x6d, 0x49, 0xec, 0x68, 0xa3, 0x8c, 0x8a, 0x4b, 0x4f, 0x61, 0x51, 0x88, 0x97, 0x16,
	0xf4, 0xe2, 0xbd, 0x07, 0x11, 0x3c, 0xad, 0x78, 0x2e, 0x6b, 0x3b, 0x34, 0x0b, 0x6a, 0xea, 0xce,
	0x54, 0xfc, 0x7c, 0xe9, 0xae, 0x09, 0xc4, 0x4b, 0x6f, 0x6f, 0xde, 0x9b, 0xe1, 0xbd, 0x37, 0x70,
	0xb6, 0x09, 0x8d, 0x34, 0xb3, 0x65, 0xed, 0x64, 0x1a, 0x21, 0x0e, 0x77, 0xd8, 0xd4, 0x30, 0x9e,
	0x07, 0x72, 0x42, 0x96, 0xbe, 0xb6, 0xc4, 0x82, 0x57, 0x90, 0x6f, 0x99, 0xc2, 0xc2, 0xaf, 0xb4,
	0x2a, 0x55, 0x35, 0xb6, 0xd9, 0x6e, 0x7c, 0x5a, 0xe1, 0x35, 0x14, 0xe2, 0xc2, 0x9a, 0x64, 0xd1,
	0xea, 0x87, 0xa5, 0xaa, 0x46, 0xf6, 0x24, 0xb1, 0xaf, 0x69, 0x4b, 0x43, 0xfe, 0x41, 0xcc, 0x6e,
	0x4d, 0x7a, 0x10, 0xe5, 0x76, 0x34, 0x00, 0x47, 0x96, 0x78, 0xd3, 0x7c, 0x32, 0x99, 0x5b, 0x28,
	0x1e, 0x49, 0x9e, 0x3d, 0xcb, 0x3e, 0x5b, 0x73, 0x03, 0xa7, 0xdd, 0x6a, 0xba, 0x46, 0x84, 0xe1,
	0xbb, 0x67, 0xd1, 0xaa, 0x1c, 0x54, 0x23, 0x1b, 0xf1, 0xdd, 0x0f, 0x1c, 0xcf, 0x6b, 0x27, 0x2f,
	0x14, 0xbe, 0xfd, 0x92, 0x70, 0x06, 0x59, 0xaa, 0x85, 0xe7, 0xd3, 0xd8, 0xb9, 0x57, 0x72, 0x52,
	0x24, 0xb2, 0xcb, 0x73, 0x80, 0x0f, 0x90, 0xff, 0xd9, 0xe0, 0x45, 0x12, 0xfb, 0x01, 0x27, 0x97,
	0xff, 0xd8, 0xf6, 0xf2, 0x2d, 0x8b, 0xef, 0xbc, 0xff, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x09, 0x58,
	0x17, 0x48, 0x62, 0x01, 0x00, 0x00,
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
