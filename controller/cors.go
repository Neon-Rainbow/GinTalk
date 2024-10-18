package controller

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type CorsConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

type Option func(*CorsConfig)

// newCorsConfig 用于创建一个新的CorsConfig实例, 并设置默认值
func newCorsConfig(options ...Option) *CorsConfig {
	defaultCfg := &CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}
	for _, option := range options {
		option(defaultCfg)
	}
	return defaultCfg
}

// WithAllowOrigins 设置允许的跨域来源
func WithAllowOrigins(origins []string) Option {
	return func(cfg *CorsConfig) {
		cfg.AllowOrigins = origins
	}
}

// WithAllowMethods 设置允许的 HTTP 方法
func WithAllowMethods(methods []string) Option {
	return func(cfg *CorsConfig) {
		cfg.AllowMethods = methods
	}
}

// WithAllowHeaders 设置允许的 HTTP 头
func WithAllowHeaders(headers []string) Option {
	return func(cfg *CorsConfig) {
		cfg.AllowHeaders = headers
	}
}

// CorsMiddleware 用于允许跨域请求
func CorsMiddleware(options ...Option) gin.HandlerFunc {
	cfg := newCorsConfig(options...)

	return func(c *gin.Context) {
		origins := strings.Join(cfg.AllowOrigins, ", ")
		methods := strings.Join(cfg.AllowMethods, ", ")
		headers := strings.Join(cfg.AllowHeaders, ", ")

		c.Writer.Header().Set("Access-Control-Allow-Origin", origins)
		c.Writer.Header().Set("Access-Control-Allow-Methods", methods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", headers)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
