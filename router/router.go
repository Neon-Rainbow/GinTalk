package router

import (
	"GinTalk/controller"
	"GinTalk/logger"
	"GinTalk/settings"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

// SetupRouter 初始化 Gin 路由
func SetupRouter(container *dig.Container) *gin.Engine {
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

	// 创建 API v1 路由组
	v1 := r.Group("/api/v1").Use(
		controller.LimitBodySizeMiddleware(),
		requestid.New(),
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

	// 使用 Dig 注入控制器，并注册相关路由
	err := container.Invoke(func(
		authController *controller.AuthHandler,
		communityController *controller.CommunityHandler,
		postController *controller.PostHandler,
		voteController *controller.VoteHandler,
		commentController *controller.CommentController,
		voteCommentController *controller.VoteCommentController,
	) {
		// 用户登录相关路由
		v1.POST("/login", authController.LoginHandler)
		v1.POST("/signup", authController.SignUpHandler)
		v1.POST("/logout", authController.LogoutHandler)
		v1.GET("/refresh_token", authController.RefreshHandler)

		// 社区相关路由
		v1.GET("/community", communityController.CommunityHandler)
		v1.GET("/community/:id", communityController.CommunityDetailHandler)

		// 帖子相关路由
		v1.POST("/post", postController.CreatePostHandler)
		v1.GET("/post", postController.GetPostListHandler)
		v1.GET("/post/community", postController.GetPostListByCommunityID)
		v1.GET("/post/:id", postController.GetPostDetailHandler)
		v1.PUT("/post", postController.UpdatePostHandler)

		// 帖子投票相关路由
		v1.POST("/vote/post", voteController.VotePostHandler)
		v1.DELETE("/vote/post", voteController.RevokeVoteHandler)
		v1.GET("/vote/post/:id", voteController.GetVoteCountHandler)
		v1.GET("/vote/post/user", voteController.MyVoteListHandler)
		v1.GET("/vote/post/list", voteController.CheckUserVotedHandler)
		v1.GET("/vote/post/batch", voteController.GetBatchPostVoteCount)
		v1.GET("/vote/post/detail", voteController.GetPostVoteDetailHandler)

		// 评论相关路由
		v1.GET("/comment/top", commentController.GetTopComments)
		v1.GET("/comment/sub", commentController.GetSubComments)
		v1.POST("/comment", commentController.CreateComment)
		v1.PUT("/comment", commentController.UpdateComment)
		v1.DELETE("/comment", commentController.DeleteComment)
		v1.GET("/comment/count", commentController.GetCommentCount)
		v1.GET("/comment/top/count", commentController.GetTopCommentCount)
		v1.GET("/comment/sub/count", commentController.GetSubCommentCount)
		v1.GET("/comment/user/count", commentController.GetCommentCountByUserID)
		v1.GET("/comment", commentController.GetCommentByCommentID)

		// 评论投票相关路由
		v1.POST("/vote/comment", voteCommentController.VoteCommentController)
		v1.DELETE("/vote/comment", voteCommentController.RemoveVoteCommentController)
		v1.GET("/vote/comment", voteCommentController.GetVoteCommentController)
		v1.GET("/vote/comment/list", voteCommentController.GetVoteCommentListController)
	})

	if err != nil {
		zap.L().Fatal("Failed to initialize controllers", zap.Error(err))
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
