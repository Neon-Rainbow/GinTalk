package dao

import (
	"GinTalk/model"
	"context"
	"gorm.io/gorm"
)

type UserDaoInterface interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUserByUsername(ctx context.Context, username string) (*model.User, error)
	FindUserByID(ctx context.Context, userID int64) (*model.User, error)
}

type UserDao struct {
	*gorm.DB
}

func NewUserDao(db *gorm.DB) UserDaoInterface {
	return &UserDao{DB: db}
}

func (ud *UserDao) CreateUser(ctx context.Context, user *model.User) error {
	sqlStr := `INSERT INTO user (user_id, username, password) VALUES (?, ?, ?)`
	return ud.WithContext(ctx).Exec(sqlStr, user.UserID, user.Username, user.Password).Error
}

func (ud *UserDao) FindUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	sqlStr := `SELECT user_id, username, password FROM user WHERE username = ? AND delete_time = 0`
	err := ud.WithContext(ctx).Raw(sqlStr, username).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ud *UserDao) FindUserByID(ctx context.Context, userID int64) (*model.User, error) {
	var user model.User
	sqlStr := `SELECT user_id, username, password FROM user WHERE user_id = ? AND delete_time = 0`
	err := ud.WithContext(ctx).Raw(sqlStr, userID).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
