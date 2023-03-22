// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/user.proto

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

// Api Endpoints for UserService service

func NewUserServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for UserService service

type UserService interface {
	TeamLogin(ctx context.Context, in *TeamLoginRequest, opts ...client.CallOption) (*LoginResponse, error)
	StudentLogin(ctx context.Context, in *StudentLoginRequest, opts ...client.CallOption) (*LoginResponse, error)
	GetInfo(ctx context.Context, in *GetInfoRequest, opts ...client.CallOption) (*UserInfoResponse, error)
	GetProfile(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*UserProfile, error)
	List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
	ListMessage(ctx context.Context, in *ListMessageRequest, opts ...client.CallOption) (*ListMessageResponse, error)
	CreateMessage(ctx context.Context, in *CreateMessageRequest, opts ...client.CallOption) (*Response, error)
	UpdateInfo(ctx context.Context, in *UpdateInfoRequest, opts ...client.CallOption) (*Response, error)
}

type userService struct {
	c    client.Client
	name string
}

func NewUserService(name string, c client.Client) UserService {
	return &userService{
		c:    c,
		name: name,
	}
}

func (c *userService) TeamLogin(ctx context.Context, in *TeamLoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.TeamLogin", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) StudentLogin(ctx context.Context, in *StudentLoginRequest, opts ...client.CallOption) (*LoginResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.StudentLogin", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) GetInfo(ctx context.Context, in *GetInfoRequest, opts ...client.CallOption) (*UserInfoResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.GetInfo", in)
	out := new(UserInfoResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) GetProfile(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*UserProfile, error) {
	req := c.c.NewRequest(c.name, "UserService.GetProfile", in)
	out := new(UserProfile)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) List(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.List", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) ListMessage(ctx context.Context, in *ListMessageRequest, opts ...client.CallOption) (*ListMessageResponse, error) {
	req := c.c.NewRequest(c.name, "UserService.ListMessage", in)
	out := new(ListMessageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) CreateMessage(ctx context.Context, in *CreateMessageRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "UserService.CreateMessage", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userService) UpdateInfo(ctx context.Context, in *UpdateInfoRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "UserService.UpdateInfo", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	TeamLogin(context.Context, *TeamLoginRequest, *LoginResponse) error
	StudentLogin(context.Context, *StudentLoginRequest, *LoginResponse) error
	GetInfo(context.Context, *GetInfoRequest, *UserInfoResponse) error
	GetProfile(context.Context, *GetRequest, *UserProfile) error
	List(context.Context, *ListRequest, *ListResponse) error
	ListMessage(context.Context, *ListMessageRequest, *ListMessageResponse) error
	CreateMessage(context.Context, *CreateMessageRequest, *Response) error
	UpdateInfo(context.Context, *UpdateInfoRequest, *Response) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) error {
	type userService interface {
		TeamLogin(ctx context.Context, in *TeamLoginRequest, out *LoginResponse) error
		StudentLogin(ctx context.Context, in *StudentLoginRequest, out *LoginResponse) error
		GetInfo(ctx context.Context, in *GetInfoRequest, out *UserInfoResponse) error
		GetProfile(ctx context.Context, in *GetRequest, out *UserProfile) error
		List(ctx context.Context, in *ListRequest, out *ListResponse) error
		ListMessage(ctx context.Context, in *ListMessageRequest, out *ListMessageResponse) error
		CreateMessage(ctx context.Context, in *CreateMessageRequest, out *Response) error
		UpdateInfo(ctx context.Context, in *UpdateInfoRequest, out *Response) error
	}
	type UserService struct {
		userService
	}
	h := &userServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&UserService{h}, opts...))
}

type userServiceHandler struct {
	UserServiceHandler
}

func (h *userServiceHandler) TeamLogin(ctx context.Context, in *TeamLoginRequest, out *LoginResponse) error {
	return h.UserServiceHandler.TeamLogin(ctx, in, out)
}

func (h *userServiceHandler) StudentLogin(ctx context.Context, in *StudentLoginRequest, out *LoginResponse) error {
	return h.UserServiceHandler.StudentLogin(ctx, in, out)
}

func (h *userServiceHandler) GetInfo(ctx context.Context, in *GetInfoRequest, out *UserInfoResponse) error {
	return h.UserServiceHandler.GetInfo(ctx, in, out)
}

func (h *userServiceHandler) GetProfile(ctx context.Context, in *GetRequest, out *UserProfile) error {
	return h.UserServiceHandler.GetProfile(ctx, in, out)
}

func (h *userServiceHandler) List(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.UserServiceHandler.List(ctx, in, out)
}

func (h *userServiceHandler) ListMessage(ctx context.Context, in *ListMessageRequest, out *ListMessageResponse) error {
	return h.UserServiceHandler.ListMessage(ctx, in, out)
}

func (h *userServiceHandler) CreateMessage(ctx context.Context, in *CreateMessageRequest, out *Response) error {
	return h.UserServiceHandler.CreateMessage(ctx, in, out)
}

func (h *userServiceHandler) UpdateInfo(ctx context.Context, in *UpdateInfoRequest, out *Response) error {
	return h.UserServiceHandler.UpdateInfo(ctx, in, out)
}
