package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// TimeoutMiddleware 请求超时中间件
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 将带有超时的上下文对象替换到原有的上下文对象中
		c.Request = c.Request.WithContext(ctx)

		// 创建一个channel用于接收请求是否处理完成
		// 使用struct{}类型是为了节省内存
		finished := make(chan struct{})

		go func() {
			c.Next() // 继续处理请求
			finished <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "请求超时"})
		case <-finished:

		}
	}
}
