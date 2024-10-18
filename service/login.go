package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/pkg"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"GinTalk/pkg/jwt"
	"context"
	"github.com/jinzhu/copier"
)

type AuthServiceInterface interface {
	LoginService(ctx context.Context, dto *DTO.LoginRequestDTO) (*DTO.LoginResponseDTO, *apiError.ApiError)
	SignupService(ctx context.Context, dto *DTO.SignUpRequestDTO) *apiError.ApiError
	RefreshTokenService(ctx context.Context, token string) (string, string, *apiError.ApiError)
}

type AuthService struct {
	userDao dao.IUserDo
}

func NewAuthService(userDao dao.IUserDo) AuthServiceInterface {
	return &AuthService{
		userDao: userDao,
	}
}

func (as *AuthService) LoginService(ctx context.Context, dto *DTO.LoginRequestDTO) (*DTO.LoginResponseDTO, *apiError.ApiError) {
	user, err := as.userDao.WithContext(ctx).
		Select(dao.User.UserID, dao.User.Username, dao.User.Password).
		Where(dao.User.Username.Eq(dto.Username)).
		First()
	if err != nil || user == nil {
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
	//err = dao.SaveUserTokenToRedis(ctx, user.ID, accessToken, refreshToken)
	//if err != nil {
	//	return nil, &apiError.ApiError{
	//		Code: code.ServerError,
	//		Msg:  "存储token失败",
	//	}
	//}
	return &DTO.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Username:     user.Username,
	}, nil
}

func (as *AuthService) SignupService(ctx context.Context, dto *DTO.SignUpRequestDTO) *apiError.ApiError {
	dto.Password = pkg.EncryptPassword(dto.Password)
	var user model.User

	err := copier.Copy(&user, dto)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "注册失败",
		}
	}

	err = as.userDao.WithContext(ctx).Create(&user)
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "注册失败",
		}
	}
	return nil
}

func (as *AuthService) RefreshTokenService(ctx context.Context, token string) (string, string, *apiError.ApiError) {
	myClaims, err := jwt.ParseToken(token)
	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.UserRefreshTokenError,
			Msg:  err.Error(),
		}
	}
	userID := myClaims.UserID
	//// 从 redis 中获取 refresh token,判断是否一致
	//value, err := dao.GetUserTokenFromRedis(ctx, userID)
	//if err != nil {
	//	return "", "", &apiError.ApiError{
	//		Code: code.UserRefreshTokenError,
	//		Msg:  err.Error(),
	//	}
	//}
	//if value["refresh_token"] != token {
	//	return "", "", &apiError.ApiError{
	//		Code: code.UserRefreshTokenError,
	//		Msg:  "refresh token 错误",
	//	}
	//}
	// 生成新的 token
	accessToken, refreshToken, err := jwt.GenerateToken(userID)
	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "生成token失败",
		}
	}
	// 更新 redis 中的 token
	//err = dao.SaveUserTokenToRedis(ctx, userID, accessToken, refreshToken)
	//if err != nil {
	//	return "", "", &apiError.ApiError{
	//		Code: code.ServerError,
	//		Msg:  "存储token失败",
	//	}
	//}
	return accessToken, refreshToken, nil

}
