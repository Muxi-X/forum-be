package service

import (
	pbf "forum-feed/proto"
	"forum/pkg/handler"

	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/spf13/viper"
	"go-micro.dev/v4/registry"

	micro "go-micro.dev/v4"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
)

var FeedClient pbf.FeedService

func FeedInit() {
	r := etcd.NewRegistry(
		registry.Addrs(viper.GetString("etcd.addr")),
		etcd.Auth(viper.GetString("etcd.username"), viper.GetString("etcd.password")),
	)
	service := micro.NewService(
		micro.Name("forum.cli.feed"),
		micro.Registry(r),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))

	service.Init()

	FeedClient = pbf.NewFeedService("forum.service.feed", service.Client())
}
