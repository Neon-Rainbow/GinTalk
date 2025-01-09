package controller

import (
	"GinTalk/DTO"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
func LoginHandler(c *gin.Context) {
	var loginDTO DTO.LoginRequestDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("LoginHandler.ShouldBindJSON() 失败", zap.Error(err))
		return
	}

	ctx := c.Request.Context()

	resp, apiError := service.LoginService(ctx, &loginDTO)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("AuthServiceInterface.LoginService() 失败", zap.Error(apiError))
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
func SignUpHandler(c *gin.Context) {
	var SignupDTO DTO.SignUpRequestDTO
	if err := c.ShouldBindJSON(&SignupDTO); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		zap.L().Error("SignUpHandler.ShouldBindJSON() 失败", zap.Error(err))
		return
	}
	ctx := c.Request.Context()

	if apiError := service.SignupService(ctx, &SignupDTO); apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("AuthServiceInterface.SignupService() 失败", zap.Error(apiError))
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
func RefreshHandler(c *gin.Context) {
	ctx := c.Request.Context()
	oldRefreshToken := c.Query("refresh_token")
	accessToken, refreshToken, apiError := service.RefreshTokenService(ctx, oldRefreshToken)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		zap.L().Error("AuthServiceInterface.RefreshTokenService() 失败", zap.Error(apiError))
		return
	}
	go func() {
		err := service.LogoutService(context.Background(), oldRefreshToken)
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
func LogoutHandler(c *gin.Context) {
	ctx := c.Request.Context()
	refreshToken := c.Query("refresh_token")
	accessToken := c.Query("access_token")
	if apiError := service.LogoutService(ctx, accessToken, refreshToken); apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, nil)
	return
}
