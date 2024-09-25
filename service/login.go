package service

import (
	"context"
	"forum-gin/DTO"
	"forum-gin/dao"
	"forum-gin/model"
	"forum-gin/pkg"
	"forum-gin/pkg/apiError"
	"forum-gin/pkg/code"
	"forum-gin/pkg/jwt"
	"github.com/jinzhu/copier"
	"time"
)

func LoginService(ctx context.Context, dto *DTO.LoginRequestDTO) (*DTO.LoginResponseDTO, *apiError.ApiError) {
	user, err := dao.GetUserByUsername(ctx, dto.Username)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.UserNotExist,
			Msg:  "用户不存在",
		}
	}
	if pkg.EncryptPassword(dto.Password) != user.Password {
		return nil, &apiError.ApiError{
			Code: code.PasswordError,
			Msg:  "密码错误",
		}
	}
	accessToken, refreshToken, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "生成token失败",
		}
	}
	err = dao.StoreKeyValue(ctx, user.ID, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, 7*24*time.Hour)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "存储token失败",
		}
	}
	return &DTO.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Username:     user.Username,
	}, nil
}

func SignupService(ctx context.Context, dto *DTO.SignUpRequestDTO) *apiError.ApiError {
	dto.Password = pkg.EncryptPassword(dto.Password)
	var user model.User

	err := copier.Copy(&user, dto)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "注册失败",
		}
	}

	err = dao.CreateUser(ctx, &user)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "注册失败",
		}
	}
	return nil
}
