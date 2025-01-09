package controller

import (
	"github.com/gin-gonic/gin"
)

// getCurrentUserID 获取当前登录用户的ID
func getCurrentUserID(c *gin.Context) (userID int64, exist bool) {
	userID = c.GetInt64(ContextUserIDKey)
	if userID == 0 {
		return 0, false
	}
	return userID, true
}

func getCurrentUsername(c *gin.Context) (username string, exist bool) {
	username = c.GetString(ContextUsernameKey)
	if username == "" {
		return "", false
	}
	return username, true
}

// isUserIDMatch 检查给定的 userID 是否与从 gin.Context 中提取的当前用户 ID 匹配。
// 如果 ID 匹配，则返回 true，否则返回 false。
//
// 参数:
//   - c: *gin.Context - 从中提取当前用户 ID 的上下文
//   - userID: int64 - 要检查的用户 ID
func isUserIDMatch(c *gin.Context, userID int64) bool {
	currentUserID, exist := getCurrentUserID(c)
	if !exist {
		return false
	}
	return currentUserID == userID
}
