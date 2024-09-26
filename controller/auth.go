package controller

import (
	"GinTalk/dao"
	"GinTalk/pkg/code"
	"GinTalk/pkg/jwt"
	"GinTalk/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

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
		if myClaims.TokenType != "access" {
			ResponseErrorWithCode(c, code.InvalidAuth)
			c.Abort()
			return
		}

		value, err := dao.GetKeyValue[map[string]string](c.Request.Context(), utils.GenerateRedisKey(utils.UserTokenKeyTemplate, myClaims.UserID))
		if err != nil {
			ResponseErrorWithCode(c, code.InvalidAuth)
			c.Abort()
			return
		}
		if value["access_token"] != token {
			ResponseErrorWithCode(c, code.InvalidAuth)
			c.Abort()
			return
		}

		c.Set("user_id", myClaims.UserID)
		c.Next()
	}
}
