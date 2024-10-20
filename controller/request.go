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

func isUserIDMatch(c *gin.Context, userID int64) bool {
	//currentUserID, err := getCurrentUserID(c)
	//if err != nil {
	//	return false
	//}
	//return currentUserID == userID
	return true // 暂时不做权限控制
}
