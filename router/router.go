package router

import (
	"GinTalk/controller"
	"GinTalk/logger"
	"GinTalk/settings"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// SetupRouter 初始化 Gin 路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 日志中间件
	r.Use(logger.GinLogger(zap.L()), logger.GinRecovery(zap.L(), true))

	// 根据配置设置 Gin 的模式
	switch settings.GetConfig().Mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 注册 prometheus 监控路由, 用于监控应用的性能
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 注册 Prometheus 中间件, 用于统计接口访问次数
	r.Use(controller.PrometheusMiddleware())

	// 创建 API v1 路由组
	v1 := r.Group("/api/v1").Use(
		controller.LimitBodySizeMiddleware(),
		requestid.New(),
		controller.TimeoutMiddleware(),
		controller.CorsMiddleware(
			controller.WithAllowOrigins([]string{"localhost"}),
		),
	)

	// 设置路由
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 用户登录相关路由
	v1.POST("/login", controller.LoginHandler)
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/logout", controller.LogoutHandler)
	v1.GET("/refresh_token", controller.RefreshHandler)

	v1.Use(controller.JWTAuthMiddleware())
	{
		// 社区相关路由
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		// 帖子相关路由
		v1.POST("/post", controller.CreatePostHandler)
		v1.DELETE("/post", controller.DeletePostHandler)
		v1.GET("/post", controller.GetPostListHandler)
		v1.GET("/post/community", controller.GetPostListByCommunityID)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.PUT("/post", controller.UpdatePostHandler)

		// 帖子投票相关路由
		v1.POST("/vote/post", controller.VotePostHandler)
		v1.DELETE("/vote/post", controller.RevokeVoteHandler)
		v1.GET("/vote/post/:id", controller.GetVoteCountHandler)
		v1.GET("/vote/post/user", controller.MyVoteListHandler)
		v1.GET("/vote/post/list", controller.CheckUserVotedHandler)
		v1.GET("/vote/post/batch", controller.GetBatchPostVoteCountHandler)
		v1.GET("/vote/post/detail", controller.GetPostVoteDetailHandler)

		// 评论相关路由
		v1.GET("/comment/top", controller.GetTopComments)
		v1.GET("/comment/sub", controller.GetSubComments)
		v1.POST("/comment", controller.CreateComment)
		v1.PUT("/comment", controller.UpdateComment)
		v1.DELETE("/comment", controller.DeleteComment)
		v1.GET("/comment/count", controller.GetCommentCount)
		v1.GET("/comment/top/count", controller.GetTopCommentCount)
		v1.GET("/comment/sub/count", controller.GetSubCommentCount)
		v1.GET("/comment/user/count", controller.GetCommentCountByUserID)
		v1.GET("/comment", controller.GetCommentByCommentID)

		// 评论投票相关路由
		v1.POST("/vote/comment", controller.VoteCommentController)
		v1.DELETE("/vote/comment", controller.RevokeVoteHandler)
		v1.GET("/vote/comment", controller.GetVoteCommentController)
		v1.GET("/vote/comment/list", controller.GetVoteCommentListController)

		v1.GET("/ws", controller.WebsocketHandle)
	}

	// 404 和 405 路由处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "请求的资源不存在",
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "请求方式非法",
		})
	})

	return r
}
