package main

import (
	"forum-post/dao"
	pb "forum-post/proto"
	"forum-post/service"
	"forum/config"
	logger "forum/log"
	"forum/pkg/handler"
	"forum/pkg/identity"
	"forum/pkg/tracer"
	"log"

	"forum/client"

	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/registry"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"

	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/spf13/viper"
)

func init() {
	// 预加载.env文件,用于本地开发
	_ = godotenv.Load()
}

func main() {
	// init config
	if err := config.Init("", "FORUM_POST"); err != nil {
		panic(err)
	}

	traceAddr := "http://" + viper.GetString("tracing.jager") + "/api/traces"
	t, io, err := tracer.NewTracer(viper.GetString("local_name"), traceAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	defer logger.SyncLogger()

	// set var t to Global Tracer (opentracing single instance mode)
	opentracing.SetGlobalTracer(t)
	r := etcd.NewRegistry(
		registry.Addrs(viper.GetString("etcd.addr")),
		etcd.Auth(viper.GetString("etcd.username"), viper.GetString("etcd.password")),
	)

	srv := micro.NewService(
		micro.Name(identity.Prefix()+viper.GetString("local_name")),
		micro.WrapHandler(
			opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapHandler(handler.ServerErrorHandlerWrapper()),

		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()),

		micro.Registry(r),
	)

	// Init will parse the command line flags.
	srv.Init()
	client.UserInit(srv)
	dao.Init()

	// Register handler
	if err := pb.RegisterPostServiceHandler(srv.Server(), service.New(dao.GetDao())); err != nil {
		panic(err)
	}

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}
