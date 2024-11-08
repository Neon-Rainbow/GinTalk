package controller

import (
	"github.com/gin-gonic/gin"
)

// getCurrentUserID 获取当前登录用户的ID
func getCurrentUserID(c *gin.Context) (userID int64, exist bool) {
	_userID, ok := c.Get(ContextUserIDKey)
	if !ok {
		return 0, false
	}
	userID, ok = _userID.(int64)
	if !ok {
		return 0, false
	}
	return
}

func getCurrentUsername(c *gin.Context) (username string, exist bool) {
	_username, ok := c.Get(ContextUsernameKey)
	if !ok {
		return "", false
	}
	username, ok = _username.(string)
	if !ok {
		return "", false
	}
	return
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
