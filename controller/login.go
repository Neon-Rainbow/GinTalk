package controller

import (
	"forum-gin/DTO"
	"forum-gin/pkg/code"
	"forum-gin/service"
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
	}

	ctx := c.Request.Context()

	resp, apiError := service.LoginService(ctx, &loginDTO)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
	}
	ResponseSuccess(c, resp)
}

func SignUpHandler(c *gin.Context) {
	var SignupDTO DTO.SignUpRequestDTO
	if err := c.ShouldBindJSON(&SignupDTO); err != nil {
		ResponseErrorWithMsg(c, code.InvalidParam, err.Error())
	}
	ctx := c.Request.Context()

	apiError := service.SignupService(ctx, &SignupDTO)
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
	}
	ResponseSuccess(c, nil)
}
