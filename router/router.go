package router

import (
	"GinTalk/controller"
	"GinTalk/logger"
	"GinTalk/settings"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// 日志中间件
	r.Use(logger.GinLogger(zap.L()), logger.GinRecovery(zap.L(), true))

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

	v1 := r.Group("/api/v1").Use(
		controller.LimitBodySizeMiddleware(),
		requestid.New(),
		//controller.TimeoutMiddleware(
		//	controller.WithTimeout(time.Duration(settings.GetConfig().Timeout)),
		//	controller.WithTimeoutMsg("请求超时"),
		//),
		controller.CorsMiddleware(
			controller.WithAllowOrigins([]string{"localhost"}),
		),
	)
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// 用户登录注册
		authController := controller.NewAuthHandler()

		v1.POST("/login", authController.LoginHandler)
		v1.POST("/signup", authController.SignUpHandler)
		v1.POST("/logout", authController.LogoutHandler)
		v1.GET("/refresh_token", authController.RefreshHandler)
	}

	//v1.Use(controller.JWTAuthMiddleware())
	communityController := controller.NewCommunityController()
	postController := controller.NewPostHandler()
	voteController := controller.NewVoteHandle()
	commentController := controller.NewCommentController()
	voteCommentController := controller.NewVoteCommentController()
	{
		v1.GET("/community", communityController.CommunityHandler)
		v1.GET("/community/:id", communityController.CommunityDetailHandler)

		v1.POST("/post", postController.CreatePostHandler)
		v1.GET("/post", postController.GetPostListHandler)
		v1.GET("/post/community", postController.GetPostListByCommunityID)
		v1.GET("/post/:id", postController.GetPostDetailHandler)
		v1.PUT("/post", postController.UpdatePostHandler)

		v1.POST("/vote/post", voteController.VotePostHandler)
		v1.DELETE("/vote/post", voteController.RevokeVoteHandler)
		v1.GET("/vote/post/:id", voteController.GetVoteCountHandler)
		v1.GET("/vote/post/user", voteController.MyVoteListHandler)
		v1.GET("/vote/post/list", voteController.CheckUserVotedHandler)
		v1.GET("/vote/post/batch", voteController.GetBatchPostVoteCount)
		v1.GET("/vote/post/detail", voteController.GetPostVoteDetailHandler)

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

		v1.POST("/vote/comment", voteCommentController.VoteCommentController)
		v1.DELETE("/vote/comment", voteCommentController.RemoveVoteCommentController)
		v1.GET("/vote/comment", voteCommentController.GetVoteCommentController)
		v1.GET("/vote/comment/list", voteCommentController.GetVoteCommentListController)
	}

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
