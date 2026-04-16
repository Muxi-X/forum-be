package main

import (
	"errors"
	"fmt"
	"forum-gateway/dao"
	"forum-gateway/router"
	"forum-gateway/router/middleware"
	forumclient "forum/client"
	"forum/config"
	"forum/log"
	"forum/pkg/handler"
	"forum/pkg/tracer"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-micro/plugins/v4/registry/kubernetes"
	opentracingWrapper "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/joho/godotenv"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go-micro.dev/v4"
	microclient "go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

func init() {
	// 预加载.env文件,用于本地开发.
	_ = godotenv.Load()
}

// @Title forum-gateway
// @Version 1.0
// @Description The gateway of forum
// @Host forum.muxixyz.com
// @BasePath /api/v1

// @tag.name user
// @tag.description 用户服务
// @tag.name auth
// @tag.description 认证服务
// @tag.name chat
// @tag.description 聊天服务
// @tag.name post
// @tag.description 帖子服务
// @tag.name feed
// @tag.description 动态服务
// @tag.name collection
// @tag.description 收藏服务
// @tag.name comment
// @tag.description 评论服务
// @tag.name like
// @tag.description 点赞服务
// @tag.name report
// @tag.description 举报服务

func main() {
	pflag.Parse()

	// init config
	if err := config.Init(*cfg, "GATEWAY"); err != nil {
		panic(err)
	}

	traceAddr := "http://" + viper.GetString("tracing.jager") + "/api/traces"
	t, io, err := tracer.NewTracer(viper.GetString("local_name"), traceAddr)
	if err != nil {
		log.Fatal("tracer error", zap.String("reason", err.Error()))
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	r := etcd.NewRegistry(
		registry.Addrs(viper.GetString("etcd.addr")),
		etcd.Auth(viper.GetString("etcd.username"), viper.GetString("etcd.password")),
	)
	service := micro.NewService(
		micro.Name("forum.gateway.client"),

		// google.protobuf.Timestamp
		// 遇到这种 proto 的时间类型，默认的 json 编码会有问题，使用 protobuf 的编码方式就不会有问题了
		micro.Client(
			microclient.NewClient(
				microclient.ContentType("application/protobuf"),
			),
		),
		micro.Registry(r),
		micro.WrapClient(
			opentracingWrapper.NewClientWrapper(opentracing.GlobalTracer()),
		),
		micro.WrapCall(handler.ClientErrorHandlerWrapper()))

	service.Init()

	// logger sync
	defer log.SyncLogger()
	forumclient.UserInit(service)
	forumclient.ChatInit(service)
	forumclient.PostInit(service)
	forumclient.FeedInit(service)
	dao.Init()
	// Set gin mode.
	gin.SetMode(viper.GetString("runmode"))

	// Create the Gin engine.
	g := gin.New()
	// Routes.
	//
	router.Load(
		// Cores.
		g,

		// MiddleWares.
		middleware.Logging(),
		ginhttp.Middleware(opentracing.GlobalTracer()),
		middleware.RequestId(),
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", zap.String("reason", err.Error()))
		}
		log.Info(fmt.Sprintf("The router has been deployed on %s successfully.", viper.GetString("addr")))
	}()

	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
