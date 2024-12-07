package controller

import (
	"GinTalk/pkg/apiError"
	"GinTalk/pkg/code"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 响应结构体
// 通过该结构体返回响应数据
// 字段:
//   - Code: 响应码
//   - Msg: 响应消息
//   - Data: 响应数据
type Response struct {
	Code code.RespCode `json:"code"`
	Msg  string        `json:"msg,omitempty"`
	Data any           `json:"data,omitempty"`
}

// ResponseSuccess 成功响应
// 返回 200 状态码
//
// 参数:
//   - c: gin.Context
//   - data: 响应数据
//
// 使用示例:
//
//	ResponseSuccess(c, data)
func ResponseSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code: code.Success,
		Msg:  "success",
		Data: data,
	})
}

// ResponseErrorWithCode 通过项目错误码来返回json数据
func ResponseErrorWithCode(c *gin.Context, respCode code.RespCode) {
	ResponseErrorWithMsg(c, respCode, respCode.GetMsg())
	return
}

func ResponseErrorWithMsg(c *gin.Context, respCode code.RespCode, msg string) {
	switch respCode {
	case code.InvalidParam:
		ResponseBadRequest(c, msg)
		return
	case code.InvalidAuth:
		ResponseUnAuthorized(c, msg)
		return
	case code.TimeOut:
		ResponseTimeout(c, msg)
		return
	case code.ServerError:
		ResponseInternalServerError(c, msg)
		return
	default:
		c.JSON(http.StatusBadRequest, Response{
			Code: respCode,
			Msg:  msg,
			Data: nil,
		})
		return
	}
}

func ResponseErrorWithApiError(c *gin.Context, apiError *apiError.ApiError) {
	switch apiError.Code {
	case code.InvalidParam:
		ResponseBadRequest(c, apiError.Msg)
		return
	case code.InvalidAuth:
		ResponseUnAuthorized(c, apiError.Msg)
		return
	case code.TimeOut:
		ResponseTimeout(c, apiError.Msg)
		return
	case code.ServerError:
		ResponseInternalServerError(c, apiError.Msg)
		return
	default:
		c.JSON(http.StatusBadRequest, Response{
			Code: apiError.Code,
			Msg:  apiError.Msg,
			Data: nil,
		},
		)
		return
	}
}

// ResponseNoContent 无内容响应
// 返回 204 状态码
func ResponseNoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// ResponseCreated 创建成功响应
// 返回 201 状态码
func ResponseCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code: code.Success,
		Msg:  "资源创建成功",
		Data: data,
	})
}

// ResponseBadRequest 参数错误响应
// 返回 400 状态码
func ResponseBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Code: code.InvalidParam,
		Msg:  msg,
		Data: nil,
	})
}

// ResponseUnAuthorized 未授权响应
// 返回 401 状态码
func ResponseUnAuthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: code.InvalidAuth,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseTimeout(c *gin.Context, msg string) {
	c.JSON(http.StatusRequestTimeout, Response{
		Code: code.TimeOut,
		Msg:  msg,
		Data: nil,
	})
}

// ResponseInternalServerError 服务器内部错误响应
// 返回 500 状态码
func ResponseInternalServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: code.ServerError,
		Msg:  msg,
		Data: nil,
	})
}
