package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type limitBodySizeConfig struct {
	LimitBytes int64
}

type limitBodySizeOption func(*limitBodySizeConfig)

// LimitBodySizeMiddleware 限制请求体大小中间件
func LimitBodySizeMiddleware(options ...limitBodySizeOption) gin.HandlerFunc {
	cfg := newLimitBodySizeConfig(options...)
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, cfg.LimitBytes)
		c.Next()
	}
}

// WithLimitBodySizeOption 设置请求体大小限制
func WithLimitBodySizeOption(limitBytes int64) limitBodySizeOption {
	return func(config *limitBodySizeConfig) {
		config.LimitBytes = limitBytes
	}
}

func newLimitBodySizeConfig(options ...limitBodySizeOption) *limitBodySizeConfig {
	config := limitBodySizeConfig{
		LimitBytes: 1024 * 1024,
	}
	for _, option := range options {
		option(&config)
	}
	return &config
}
