package controller

import (
	"GinTalk/pkg/code"
	"GinTalk/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	// ContextUserIDKey 是上下文中用户ID的key
	ContextUserIDKey   = "user_id"
	ContextUsernameKey = "username"
)

// JWTAuthMiddleware JWT 认证中间件, 用于验证用户是否登录
// 如果用户登录, 会将用户ID设置到上下文中
// 如果用户未登录, 会返回错误响应
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			ResponseErrorWithCode(c, code.InvalidAuth)
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ResponseErrorWithCode(c, code.InvalidAuth)
			c.Abort()
			return
		}
		token := parts[1]
		myClaims, err := jwt.ParseToken(token)
		if err != nil {
			ResponseErrorWithCode(c, code.InvalidAuth)
			c.Abort()
			return
		}
		if myClaims.TokenType != jwt.AccessTokenName {
			ResponseErrorWithCode(c, code.InvalidAuth)
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, myClaims.UserID)
		c.Set(ContextUsernameKey, myClaims.Username)
		c.Next()
	}
}
