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
// 登录服务，根据用户名和密码查询用户，如果用户存在且密码正确，则生成token返回
//
// 参数
//   - ctx: 上下文
//   - dto: 登录请求数据传输对象
//
// 返回值
//   - *DTO.LoginResponseDTO: 登录响应数据传输对象
//   - *apiError.ApiError: 错误信息
//
// 使用示例
//
//	resp, apiError := service.LoginService(ctx, &loginDTO)
//	if apiError != nil {
//	  ResponseErrorWithApiError(c, apiError)
//	  zap.L().Error("AuthServiceInterface.LoginService() 失败", zap.Error(apiError))
//	  return
//	}
//	ResponseSuccess(c, resp)
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
// SignupService 处理用户注册过程。
// 它接受一个上下文和一个 SignUpRequestDTO 作为输入，加密密码，
// 将 DTO 复制到 User 模型，生成唯一的用户 ID，并在数据库中创建用户。
// 如果任何步骤失败，它将返回一个包含适当错误代码和消息的 ApiError。
//
// 参数:
//   - ctx: 用于管理请求范围值、取消和截止日期的上下文。
//   - dto: 包含用户注册详细信息的 SignUpRequestDTO 指针。
//
// 返回值:
//   - *apiError.ApiError: 如果任何步骤失败，则返回包含错误代码和消息的错误对象，否则返回 nil。
//
// 使用示例:
//
//	apiError := service.SignupService(ctx, &signupDTO)
//	if apiError != nil {
//	  ResponseErrorWithApiError(c, apiError)
//	  zap.L().Error("AuthServiceInterface.SignupService() 失败", zap.Error(apiError))
//	  return
//	}
//	ResponseSuccess(c, nil)
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
// RefreshTokenService 从传入的 token 中解析出用户 ID 和用户名，然后生成新的访问令牌和刷新令牌。
// 如果解析失败，它将返回一个包含错误代码和消息的 ApiError。
//
// 参数:
//   - ctx: 用于管理请求范围值、取消和截止日期的上下文。
//   - token: 包含用户 ID 和用户名的 token 字符串。
//
// 返回值:
//   - string: 新的访问令牌。
//   - string: 新的刷新令牌。
//   - *apiError.ApiError: 如果解析失败，则返回包含错误代码和消息的错误对象，否则返回 nil。
//
// 使用示例:
//
//	accessToken, refreshToken, apiError := service.RefreshTokenService(ctx, oldRefreshToken)
//	if apiError != nil {
//	  ResponseErrorWithApiError(c, apiError)
//	  zap.L().Error("AuthServiceInterface.RefreshTokenService() 失败", zap.Error(apiError))
//	  return
//	}
//	go func() {
//	  err := service.LogoutService(context.Background(), oldRefreshToken)
//	  if err != nil {
//	    zap.L().Error("refresh token logout failed", zap.Error(err))
//	  }
//	}()
//	ResponseSuccess(c, gin.H{
//	  "access_token":  accessToken,
//	  "refresh_token": refreshToken,
//	})
//	return
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

// LogoutService 退出登录
// LogoutService 将传入的 token 添加到黑名单中，以便用户无法再使用它。
// 如果添加失败，它将返回一个包含错误代码和消息的 ApiError。
//
// 参数:
//   - ctx: 用于管理请求范围值、取消和截止日期的上下文。
//   - token: 包含用户 ID 和用户名的 token 字符串。
//
// 返回值:
//   - *apiError.ApiError: 如果添加失败，则返回包含错误代码和消息的错误对象，否则返回 nil。
//
// 使用示例:
//
//	apiError := service.LogoutService(ctx, accessToken, refreshToken)
//	if apiError != nil {
//	  ResponseErrorWithApiError(c, apiError)
//	  return
//	}
//	ResponseSuccess(c, nil)
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
