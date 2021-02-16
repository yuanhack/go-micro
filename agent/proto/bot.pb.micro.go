// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: agent/proto/bot.proto

package go_micro_bot

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/macheal/go-micro/v2/api"
	client "github.com/macheal/go-micro/v2/client"
	server "github.com/macheal/go-micro/v2/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Command service

func NewCommandEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Command service

type CommandService interface {
	Help(ctx context.Context, in *HelpRequest, opts ...client.CallOption) (*HelpResponse, error)
	Exec(ctx context.Context, in *ExecRequest, opts ...client.CallOption) (*ExecResponse, error)
}

type commandService struct {
	c    client.Client
	name string
}

func NewCommandService(name string, c client.Client) CommandService {
	return &commandService{
		c:    c,
		name: name,
	}
}

func (c *commandService) Help(ctx context.Context, in *HelpRequest, opts ...client.CallOption) (*HelpResponse, error) {
	req := c.c.NewRequest(c.name, "Command.Help", in)
	out := new(HelpResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commandService) Exec(ctx context.Context, in *ExecRequest, opts ...client.CallOption) (*ExecResponse, error) {
	req := c.c.NewRequest(c.name, "Command.Exec", in)
	out := new(ExecResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Command service

type CommandHandler interface {
	Help(context.Context, *HelpRequest, *HelpResponse) error
	Exec(context.Context, *ExecRequest, *ExecResponse) error
}

func RegisterCommandHandler(s server.Server, hdlr CommandHandler, opts ...server.HandlerOption) error {
	type command interface {
		Help(ctx context.Context, in *HelpRequest, out *HelpResponse) error
		Exec(ctx context.Context, in *ExecRequest, out *ExecResponse) error
	}
	type Command struct {
		command
	}
	h := &commandHandler{hdlr}
	return s.Handle(s.NewHandler(&Command{h}, opts...))
}

type commandHandler struct {
	CommandHandler
}

func (h *commandHandler) Help(ctx context.Context, in *HelpRequest, out *HelpResponse) error {
	return h.CommandHandler.Help(ctx, in, out)
}

func (h *commandHandler) Exec(ctx context.Context, in *ExecRequest, out *ExecResponse) error {
	return h.CommandHandler.Exec(ctx, in, out)
}
