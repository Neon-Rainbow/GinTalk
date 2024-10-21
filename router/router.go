package router

import (
	"GinTalk/controller"
	"GinTalk/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// 日志中间件
	r.Use(logger.GinLogger(zap.L())).Use(logger.GinRecovery(zap.L(), true))

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
	{
		v1.GET("/community", communityController.CommunityHandler)
		v1.GET("/community/:id", communityController.CommunityDetailHandler)

		v1.POST("/post", postController.CreatePostHandler)
		v1.GET("/post", postController.GetPostListHandler)
		v1.GET("/post/community", postController.GetPostListByCommunityID)
		v1.GET("/post/:id", postController.GetPostDetailHandler)
		v1.PUT("/post", postController.UpdatePostHandler)

		v1.POST("/vote", voteController.VoteHandler)
		v1.DELETE("/vote", voteController.RevokeVoteHandler)
		v1.GET("/vote/:id", voteController.GetVoteCountHandler)
		v1.GET("/vote/user", voteController.MyVoteListHandler)
		v1.GET("/vote/list", voteController.CheckUserVotedHandler)
		v1.GET("/vote/post/detail", voteController.GetPostVoteDetailHandler)
		v1.GET("/vote/comment/detail", voteController.GetCommentVoteDetailHandler)

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
