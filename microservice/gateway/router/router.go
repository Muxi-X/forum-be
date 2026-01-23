package router

import (
	"forum-gateway/dao"
	_ "forum-gateway/docs"
	"forum-gateway/handler"
	"forum-gateway/handler/chat"
	"forum-gateway/handler/collection"
	"forum-gateway/handler/comment"
	"forum-gateway/handler/feed"
	"forum-gateway/handler/like"
	"forum-gateway/handler/post"
	"forum-gateway/handler/report"
	"forum-gateway/handler/sd"
	"forum-gateway/handler/user"
	"forum-gateway/router/middleware"
	"forum/pkg/constvar"
	"forum/pkg/errno"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)

	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		handler.SendError(c, errno.ErrIncorrectAPIRoute, nil, "", "")
	})

	// swagger API doc
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 权限要求，普通用户/管理员/超管
	normalRequired := middleware.AuthMiddleware(constvar.AuthLevelNormal)
	adminRequired := middleware.AuthMiddleware(constvar.AuthLevelAdmin)
	// superAdminRequired := middleware.AuthMiddleware(constvar.AuthLevelSuperAdmin)

	// auth 模块
	authRouter := g.Group("api/v1/auth/login")
	{
		authRouter.POST("/student", user.StudentLogin)
		authRouter.POST("/team", user.TeamLogin)
	}

	// user 模块
	userRouter := g.Group("api/v1/user")
	userRouter.Use(normalRequired)
	{
		userRouter.GET("/profile/:id", user.GetProfile)
		userRouter.GET("/myprofile", user.GetMyProfile)
		userRouter.GET("/list", user.List)
		userRouter.PUT("", user.UpdateInfo)
		userRouter.GET("/message/list", user.ListMessage)
		userRouter.POST("/private_message", user.CreatePrivateMessage)
		userRouter.POST("/message", adminRequired, user.CreateMessage)
		userRouter.DELETE("/private_message", user.DeletePrivateMessage)
		userRouter.GET("/private_message/list", user.ListPrivateMessage)
	}

	chatRouter := g.Group("api/v1/chat")
	chatRouter.Use(normalRequired)
	{
		chatRouter.GET("/history/:id", chat.ListHistory)
		chatRouter.GET("/ws", chat.WsHandler)
		chatRouter.GET("/userList", chat.UserList)
	}

	postRouter := g.Group("api/v1/post")
	postApi := post.New(dao.GetDao())
	postRouter.Use(normalRequired)
	{
		postRouter.GET("/list/:domain", postApi.ListMainPost)
		postRouter.GET("/published/:user_id", postApi.ListUserPost)
		postRouter.GET("/:post_id", postApi.Get)
		postRouter.POST("", postApi.Create)
		postRouter.DELETE("/:post_id", postApi.Delete)
		postRouter.PUT("", postApi.UpdateInfo)
		postRouter.GET("/popular_tag", postApi.ListPopularTag)
		postRouter.GET("/qiniu_token", postApi.GetQiNiuToken)
	}

	commentRouter := g.Group("api/v1/comment")
	commentApi := comment.New(dao.GetDao())
	commentRouter.Use(normalRequired)
	{
		commentRouter.GET("/:comment_id", commentApi.Get)
		commentRouter.POST("", commentApi.Create)
		commentRouter.DELETE("/:comment_id", commentApi.Delete)
	}

	likeRouter := g.Group("api/v1/like")
	likeApi := like.New(dao.GetDao())
	likeRouter.Use(normalRequired)
	{
		likeRouter.GET("/list/:user_id", likeApi.GetUserLikeList)
		likeRouter.POST("", likeApi.CreateOrRemove)
	}

	// feed
	feedRouter := g.Group("api/v1/feed")
	feedApi := feed.New(dao.GetDao())
	feedRouter.Use(normalRequired)
	{
		feedRouter.GET("/list/:user_id", feedApi.List)
	}

	// collection
	collectionRouter := g.Group("api/v1/collection")
	collectionApi := collection.New(dao.GetDao())
	collectionRouter.Use(normalRequired)
	{
		collectionRouter.GET("/list/:user_id", collectionApi.List)
		collectionRouter.POST("/:post_id", collectionApi.CreateOrRemove)
	}

	// report
	reportRouter := g.Group("api/v1/report")
	reportApi := report.New(dao.GetDao())
	{
		reportRouter.POST("", normalRequired, reportApi.Create)
		reportRouter.GET("/list", adminRequired, reportApi.List)
		reportRouter.PUT("", adminRequired, reportApi.Handle)
	}

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
