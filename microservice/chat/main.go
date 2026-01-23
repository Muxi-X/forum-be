package main

import (
	"forum-chat/dao"
	pb "forum-chat/proto"
	"forum-chat/service"
	"forum/config"
	logger "forum/log"
	"forum/pkg/handler"
	"forum/pkg/identity"
	"forum/pkg/tracer"
	"log"

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
	// 预加载.env文件,用于本地开发.
	_ = godotenv.Load()
}

func main() {
	// init config
	if err := config.Init("", "FORUM_CHAT"); err != nil {
		panic(err)
	}

	t, io, err := tracer.NewTracer(viper.GetString("local_name"), viper.GetString("tracing.jager"))
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
		micro.Registry(r),
	)

	// Init will parse the command line flags.
	srv.Init()

	dao.Init()

	// Register handler
	if err := pb.RegisterChatServiceHandler(srv.Server(), service.New(dao.GetDao())); err != nil {
		panic(err)
	}

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}
