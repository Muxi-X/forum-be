package main

import (
	"flag"
	"forum/pkg/identity"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"

	"forum-feed/dao"
	pb "forum-feed/proto"
	"forum-feed/service"
	"forum/config"
	logger "forum/log"
	"forum/pkg/handler"
	"forum/pkg/tracer"

	_ "github.com/go-micro/plugins/v4/registry/kubernetes"

	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/spf13/viper"
	cli "github.com/urfave/cli/v2"
)

func init() {
	// 预加载.env文件,用于本地开发
	_ = godotenv.Load()
}

// 使用--sub运行subscribe服务
// 否则默认运行feed服务
var subFg = flag.Bool("sub", false, "use subscribe service mode")

// InitRpcClient ... 记录 bug: go-micro 框架的命令行参数和 go 标准库的 flag 冲突了
func initRpcClient() {
	service.UserInit()
}

// 包含两个服务：feed服务和subscribe服务
// subscribe服务 --> 异步将feed数据写入数据库
func main() {
	flag.Parse()

	var err error

	// init config
	if !*subFg {
		// feed-service
		err = config.Init("./conf/config.yaml", "FORUM_FEED")
		initRpcClient()
	} else {
		// sub-service
		err = config.Init("./conf/config_sub.yaml", "FORUM_SUB")
	}

	if err != nil {
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

	dao.Init()

	if *subFg {
		// sub-service
		logger.Info("Subscribe service start...")
		service.SubServiceRun()
		return
	}
	r := etcd.NewRegistry(
		registry.Addrs(viper.GetString("etcd.addr")),
		etcd.Auth(viper.GetString("etcd.username"), viper.GetString("etcd.password")),
	)
	// feed-service
	srv := micro.NewService(
		micro.Name("forum/"+identity.GetIdentity()+"/"+viper.GetString("local_name")),
		micro.WrapHandler(
			opentracingWrapper.NewHandlerWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapHandler(handler.ServerErrorHandlerWrapper()),
		micro.Flags(&cli.BoolFlag{
			Name:   "sub",
			Usage:  "use subscribe service mode",
			Hidden: false,
		}),
		micro.Registry(r),
	)

	// Init will parse the command line flags.
	srv.Init()

	if err := pb.RegisterFeedServiceHandler(srv.Server(), service.New(dao.GetDao())); err != nil {
		panic(err)
	}

	// Run the server
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
	}
}
