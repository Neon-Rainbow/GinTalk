package controller

import (
	"GinTalk/DTO"
	"GinTalk/container"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service.AuthServiceInterface
}

// NewAuthHandler 创建 AuthHandler 实例
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		container.GetAuthService(),
	}
}

// LoginHandler 登录接口
// @Summary 登录接口
// @Description 登录接口
// @Tags 登录
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
func (ah *AuthHandler) RefreshHandler(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Query("refresh_token")
	accessToken, refreshToken, apiError := ah.AuthServiceInterface.RefreshTokenService(ctx, token)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	return
}
