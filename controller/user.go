package controller

import (
	"GinTalk/DTO"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service.AuthServiceInterface
}

// NewAuthHandler 创建 AuthHandler 实例
func NewAuthHandler(service service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		AuthServiceInterface: service,
	}
}

// LoginHandler 登录接口
// @Summary 登录接口
// @Description 登录接口
// @Tags 登录
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} Response
// @Router /api/v1/login [post]
func (ah *AuthHandler) LoginHandler(c *gin.Context) {
	var loginDTO DTO.LoginRequestDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	ctx := c.Request.Context()

	resp, apiError := ah.AuthServiceInterface.LoginService(ctx, &loginDTO)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, resp)
	return
}

// SignUpHandler 注册接口
// @Summary 注册接口
// @Description 注册接口
// @Tags 登录
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Param email body string true "邮箱"
// @Param gender body string true "性别"
// @Success 200 {object} Response
// @Router /api/v1/signup [post]
func (ah *AuthHandler) SignUpHandler(c *gin.Context) {
	var SignupDTO DTO.SignUpRequestDTO
	if err := c.ShouldBindJSON(&SignupDTO); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	ctx := c.Request.Context()

	apiError := ah.AuthServiceInterface.SignupService(ctx, &SignupDTO)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
}

// RefreshHandler 刷新token
// @Summary 刷新token
// @Description 刷新token
// @Tags 登录
// @Accept json
// @Produce json
// @Param refresh_token query string true
// @Success 200 {object} Response
// @Router /api/v1/refresh_token [get]
func (ah *AuthHandler) RefreshHandler(c *gin.Context) {
	ctx := c.Request.Context()
	oldRefreshToken := c.Query("refresh_token")
	accessToken, refreshToken, apiError := ah.AuthServiceInterface.RefreshTokenService(ctx, oldRefreshToken)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	go func() {
		err := ah.AuthServiceInterface.LogoutService(context.Background(), oldRefreshToken)
		if err != nil {
			zap.L().Error("refresh token logout failed", zap.Error(err))
		}
	}()
	ResponseSuccess(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	return
}

// LogoutHandler 退出登录
// @Summary 退出登录
// @Description 退出登录
// @Tags 登录
// @Accept json
// @Produce json
// @Param access_token query string true
// @Param refresh_token query string true
// @Success 200 {object} Response
// @Router /api/v1/logout [post]
func (ah *AuthHandler) LogoutHandler(c *gin.Context) {
	ctx := c.Request.Context()
	refreshToken := c.Query("refresh_token")
	accessToken := c.Query("access_token")
	apiError := ah.AuthServiceInterface.LogoutService(ctx, accessToken, refreshToken)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
	return
}
