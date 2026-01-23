package client

import (
	pbc "forum-chat/proto"
	"forum/pkg/handler"

	"forum/pkg/identity"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/spf13/viper"
	"go-micro.dev/v4/registry"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
)

var ChatClient pbc.ChatService

func ChatInit() {
	r := etcd.NewRegistry(
		registry.Addrs(viper.GetString("etcd.addr")),
		etcd.Auth(viper.GetString("etcd.username"), viper.GetString("etcd.password")),
	)
	service := micro.NewService(
		micro.Name("forum.cli.chat"),
		micro.Registry(r),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))

	service.Init()

	ChatClient = pbc.NewChatService(identity.Prefix()+"forum.service.chat", service.Client())
}
