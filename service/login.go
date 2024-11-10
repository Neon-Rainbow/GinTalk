package service

import (
	"GinTalk/DTO"
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/pkg"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"GinTalk/pkg/jwt"
	"GinTalk/pkg/snowflake"
	"context"
	"time"

	"github.com/jinzhu/copier"
)

// LoginService 登录服务
func LoginService(ctx context.Context, dto *DTO.LoginRequestDTO) (*DTO.LoginResponseDTO, *apiError.ApiError) {
	user, err := dao.FindUserByUsername(ctx, dto.Username)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "登录失败",
		}
	}
	if user == nil {
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
	accessToken, refreshToken, err := jwt.GenerateToken(user.UserID, user.Username)
	if err != nil {
		return nil, &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "生成token失败",
		}
	}

	return &DTO.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.UserID,
		Username:     user.Username,
	}, nil
}

// SignupService 注册服务
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

	user.UserID, err = snowflake.GetID()
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

// RefreshTokenService 刷新token
func RefreshTokenService(ctx context.Context, token string) (string, string, *apiError.ApiError) {
	myClaims, err := jwt.ParseToken(token)
	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.UserRefreshTokenError,
			Msg:  err.Error(),
		}
	}
	if myClaims.TokenType != jwt.RefreshTokenName {
		return "", "", &apiError.ApiError{
			Code: code.UserRefreshTokenError,
			Msg:  "token类型错误",
		}
	}

	accessToken, refreshToken, err := jwt.GenerateToken(myClaims.UserID, myClaims.Username)
	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "生成token失败",
		}
	}

	err = cache.AddTokenToBlacklist(ctx, token, time.Until(myClaims.ExpiresAt.Time))

	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "刷新token失败",
		}
	}

	return accessToken, refreshToken, nil

}

func LogoutService(ctx context.Context, token ...string) *apiError.ApiError {
	for _, t := range token {
		myClaims, err := jwt.ParseToken(t)
		if err != nil {
			return &apiError.ApiError{
				Code: code.UserRefreshTokenError,
				Msg:  err.Error(),
			}
		}

		err = cache.AddTokenToBlacklist(ctx, t, time.Until(myClaims.ExpiresAt.Time))
		if err != nil {
			return &apiError.ApiError{
				Code: code.ServerError,
				Msg:  "登出失败",
			}
		}
	}

	return nil
}
