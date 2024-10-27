package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

type TimeoutConfig struct {
	Timeout         time.Duration
	TimeoutMsg      string
	ErrorHandleFunc gin.HandlerFunc
}

type TimeoutOption func(*TimeoutConfig)

func newTimeoutConfig(options ...TimeoutOption) *TimeoutConfig {
	defaultCfg := &TimeoutConfig{
		Timeout:    time.Second * 100,
		TimeoutMsg: "请求超时",
		ErrorHandleFunc: func(c *gin.Context) {
			ResponseTimeout(c, "请求超时")
			c.Abort()
		},
	}
	for _, option := range options {
		option(defaultCfg)
	}
	return defaultCfg
}

func WithTimeout(timeout time.Duration) TimeoutOption {
	return func(cfg *TimeoutConfig) {
		cfg.Timeout = timeout * time.Second
	}
}

func WithTimeoutMsg(msg string) TimeoutOption {
	return func(cfg *TimeoutConfig) {
		cfg.TimeoutMsg = msg
	}
}

func WithErrorHandle(handle gin.HandlerFunc) TimeoutOption {
	return func(cfg *TimeoutConfig) {
		cfg.ErrorHandleFunc = handle
	}
}

// TimeoutMiddleware 请求超时中间件
func TimeoutMiddleware(options ...TimeoutOption) gin.HandlerFunc {
	cfg := newTimeoutConfig(options...)
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.Timeout)
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
			cfg.ErrorHandleFunc(c)
		case <-finished:

		}
	}
}
