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
