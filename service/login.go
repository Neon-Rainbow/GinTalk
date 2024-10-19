package service

import (
	"GinTalk/DTO"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/pkg"
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"GinTalk/pkg/jwt"
	"GinTalk/pkg/snowflake"
	"context"
	"github.com/jinzhu/copier"
)

var _ AuthServiceInterface = (*AuthService)(nil)

type AuthServiceInterface interface {
	LoginService(ctx context.Context, dto *DTO.LoginRequestDTO) (*DTO.LoginResponseDTO, *apiError.ApiError)
	SignupService(ctx context.Context, dto *DTO.SignUpRequestDTO) *apiError.ApiError
	RefreshTokenService(ctx context.Context, token string) (string, string, *apiError.ApiError)
}

type AuthService struct {
	dao.IUserDo
	dao.UserDaoInterface
}

func NewAuthService(userDao dao.IUserDo, userDaoInterface dao.UserDaoInterface) AuthServiceInterface {
	return &AuthService{
		userDao,
		userDaoInterface,
	}
}

func (as *AuthService) LoginService(ctx context.Context, dto *DTO.LoginRequestDTO) (*DTO.LoginResponseDTO, *apiError.ApiError) {
	user, err := as.UserDaoInterface.FindUserByUsername(ctx, dto.Username)
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
	accessToken, refreshToken, err := jwt.GenerateToken(user.UserID)
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

	user.UserID, err = snowflake.GetID()
	if err != nil {
		return &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "注册失败",
		}
	}

	err = as.UserDaoInterface.CreateUser(ctx, &user)
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

	accessToken, refreshToken, err := jwt.GenerateToken(userID)
	if err != nil {
		return "", "", &apiError.ApiError{
			Code: code.ServerError,
			Msg:  "生成token失败",
		}
	}

	return accessToken, refreshToken, nil

}
