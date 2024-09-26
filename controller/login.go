package controller

import (
	"GinTalk/DTO"
	"GinTalk/pkg/code"
	"GinTalk/service"
	"github.com/gin-gonic/gin"
)

// LoginHandler 登录接口
// @Summary 登录接口
// @Description 登录接口
// @Tags 登录
func LoginHandler(c *gin.Context) {
	var loginDTO DTO.LoginRequestDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}

	ctx := c.Request.Context()

	resp, apiError := service.LoginService(ctx, &loginDTO)
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
func SignUpHandler(c *gin.Context) {
	var SignupDTO DTO.SignUpRequestDTO
	if err := c.ShouldBindJSON(&SignupDTO); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
		return
	}
	ctx := c.Request.Context()

	apiError := service.SignupService(ctx, &SignupDTO)
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
func RefreshHandler(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Query("refresh_token")
	accessToken, refreshToken, apiError := service.RefreshTokenService(ctx, token)
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
