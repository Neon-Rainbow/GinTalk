package controller

import (
	"GinTalk/cache"
	"GinTalk/pkg/code"
	"GinTalk/pkg/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

const (
	// ContextUserIDKey 是上下文中用户ID的key
	ContextUserIDKey = "user_id"
	// ContextUsernameKey 是上下文中用户名的key
	ContextUsernameKey = "username"
)

// JWTAuthMiddleware JWT 认证中间件, 用于验证用户是否登录
// 如果用户登录, 会将用户ID设置到上下文中
// 如果用户未登录, 会返回错误响应
// 如果 token 在黑名单中, 会返回错误响应
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			ResponseUnAuthorized(c, "请求未携带 token")
			zap.L().Info("请求未携带 token")
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ResponseUnAuthorized(c, "token 格式错误, 应为 `Bearer {token}`")
			zap.L().Info("token 格式错误")
			c.Abort()
			return
		}
		token := parts[1]
		myClaims, err := jwt.ParseToken(token)
		if err != nil {
			ResponseUnAuthorized(c, "token 解析失败")
			zap.L().Error("token 解析失败", zap.Error(err))
			c.Abort()
			return
		}
		if myClaims.TokenType != jwt.AccessTokenName {
			ResponseUnAuthorized(c, "token 类型错误")
			zap.L().Info("token 类型错误")
			c.Abort()
			return
		}

		exist, err := cache.IsTokenInBlacklist(c.Request.Context(), token)
		if exist {
			ResponseUnAuthorized(c, "token 已失效")
			zap.L().Info("token 已失效")
			c.Abort()
		}
		if err != nil {
			ResponseErrorWithMsg(c, code.ServerError, fmt.Sprintf("authCache.IsTokenInBlacklist() 出错: %v", err))
			zap.L().Error("authCache.IsTokenInBlacklist() 出错", zap.Error(err))
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, myClaims.UserID)
		c.Set(ContextUsernameKey, myClaims.Username)
		c.Next()
		return
	}
}
