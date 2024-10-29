package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// getCurrentUserID 获取当前登录用户的ID
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	_userID, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = errors.New("当前用户未登录")
		return
	}
	userID, ok = _userID.(int64)
	if !ok {
		err = errors.New("当前用户未登录")
		return
	}
	return
}

func getCurrentUsername(c *gin.Context) (username string, err error) {
	_username, ok := c.Get(ContextUsernameKey)
	if !ok {
		err = errors.New("当前用户未登录")
		return
	}
	username, ok = _username.(string)
	if !ok {
		err = errors.New("当前用户未登录")
		return
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
	currentUserID, err := getCurrentUserID(c)
	if err != nil {
		return false
	}
	return currentUserID == userID
}
