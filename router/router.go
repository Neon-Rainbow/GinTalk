package router

import (
	"GinTalk/controller"
	"GinTalk/logger"
	"GinTalk/settings"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	// 日志中间件
	r.Use(logger.GinLogger(zap.L())).Use(logger.GinRecovery(zap.L(), true))

	v1 := r.Group("/api/v1").Use(
		controller.LimitBodySizeMiddleware(1<<20),
		requestid.New(),
		controller.TimeoutMiddleware(time.Duration(settings.Conf.Timeout)*time.Second),
		controller.CorsMiddleware(controller.NewCorsConfig()),
	)
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		v1.POST("/login", controller.LoginHandler)
		v1.POST("/signup", controller.SignUpHandler)
		v1.GET("/refresh_token", controller.RefreshHandler)
	}

	v1.Use(controller.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
