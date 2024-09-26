package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/pkg"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"GinTalk/pkg/jwt"
	"GinTalk/utils"
	"context"
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
	err = dao.StoreKeyValue(ctx, utils.GenerateRedisKey(utils.UserTokenKeyTemplate, user.ID), map[string]interface{}{
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

func RefreshTokenService(ctx context.Context, token string) (string, string, *apiError.ApiError) {
	myClaims, err := jwt.ParseToken(token)
	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.UserRefreshTokenError,
			Msg:  err.Error(),
		}
	}
	userID := myClaims.UserID
	// 从 redis 中获取 refresh token,判断是否一致
	value, err := dao.GetKeyValue[map[string]string](ctx, utils.GenerateRedisKey(utils.UserTokenKeyTemplate, userID))
	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.UserRefreshTokenError,
			Msg:  err.Error(),
		}
	}
	if value["refresh_token"] != token {
		return "", "", &apiError.ApiError{
			Code: code.UserRefreshTokenError,
			Msg:  "refresh token 错误",
		}
	}
	// 生成新的 token
	accessToken, refreshToken, err := jwt.GenerateToken(userID)
	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "生成token失败",
		}
	}
	// 更新 redis 中的 token
	err = dao.StoreKeyValue(ctx, utils.GenerateRedisKey(utils.UserTokenKeyTemplate, userID), map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, 7*24*time.Hour)
	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "存储token失败",
		}
	}
	return accessToken, refreshToken, nil
}
