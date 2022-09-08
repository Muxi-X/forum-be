// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/chat.proto

package proto

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for ChatService service

func NewChatServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for ChatService service

type ChatService interface {
	Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*Response, error)
	GetList(ctx context.Context, in *GetListRequest, opts ...client.CallOption) (*GetListResponse, error)
}

type chatService struct {
	c    client.Client
	name string
}

func NewChatService(name string, c client.Client) ChatService {
	return &chatService{
		c:    c,
		name: name,
	}
}

func (c *chatService) Create(ctx context.Context, in *CreateRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "ChatService.Create", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatService) GetList(ctx context.Context, in *GetListRequest, opts ...client.CallOption) (*GetListResponse, error) {
	req := c.c.NewRequest(c.name, "ChatService.GetList", in)
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

func RegisterChatServiceHandler(s server.Server, hdlr ChatServiceHandler, opts ...server.HandlerOption) error {
	type chatService interface {
		Create(ctx context.Context, in *CreateRequest, out *Response) error
		GetList(ctx context.Context, in *GetListRequest, out *GetListResponse) error
	}
	type ChatService struct {
		chatService
	}
	h := &chatServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&ChatService{h}, opts...))
}

type chatServiceHandler struct {
	ChatServiceHandler
}

func (h *chatServiceHandler) Create(ctx context.Context, in *CreateRequest, out *Response) error {
	return h.ChatServiceHandler.Create(ctx, in, out)
}

func (h *chatServiceHandler) GetList(ctx context.Context, in *GetListRequest, out *GetListResponse) error {
	return h.ChatServiceHandler.GetList(ctx, in, out)
}
