package controller

import (
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code code.RespCode `json:"code"`
	Msg  string        `json:"msg"`
	Data interface{}   `json:"data"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code.Success,
		Msg:  "success",
		Data: data,
	})
}

func ResponseErrorWithCode(c *gin.Context, code code.RespCode) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  code.GetMsg(),
		Data: nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code code.RespCode, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseErrorWithApiError(c *gin.Context, apiError *apiError.ApiError) {
	c.JSON(http.StatusOK, Response{
		Code: apiError.Code,
		Msg:  apiError.Msg,
		Data: nil,
	})
}
