package dao

import (
	"context"
	"forum-gin/dao/MySQL"
	"forum-gin/model"
)

// GetUserByUsername 根据用户名获取用户信息
func GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := MySQL.GetDB().WithContext(ctx).Where("username = ?", username).First(&user).Error
	return &user, err
}

func CreateUser(ctx context.Context, user *model.User) error {
	err := MySQL.GetDB().WithContext(ctx).Create(user).Error
	return err
}
