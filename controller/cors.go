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

// NewCorsConfig 用于创建一个新的CorsConfig实例, 并设置默认值
func NewCorsConfig() *CorsConfig {
	return &CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}
}

// CorsMiddleware 用于允许跨域请求
func CorsMiddleware(cfg *CorsConfig) gin.HandlerFunc {
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
