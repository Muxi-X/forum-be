package main

import (
	"forum-user/dao"
	"forum-user/pkg/auth"
	"forum/pkg/identity"

	pb "forum-user/proto"
	"forum-user/service"
	"forum/config"
	logger "forum/log"
	"forum/pkg/handler"
	"forum/pkg/tracer"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"

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
	if err := config.Init("", "FORUM_USER"); err != nil {
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

	// init oauth-manager and some variables
	auth.InitVar()
	auth.OauthManager.Init()
	r := etcd.NewRegistry(
		registry.Addrs(viper.GetString("etcd.addr")),
		etcd.Auth(viper.GetString("etcd.username"), viper.GetString("etcd.password")),
	)
	srv := micro.NewService(
		micro.Name("forum/"+identity.GetIdentity()+"/"+viper.GetString("local_name")),
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
	if err := pb.RegisterUserServiceHandler(srv.Server(), service.New(dao.GetDao())); err != nil {
		panic(err)
	}

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}
