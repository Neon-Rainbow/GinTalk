package controller

import (
	"GinTalk/pkg/code"
	"GinTalk/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CommunityHandler(c *gin.Context) {
	list, apiError := service.GetCommunityList(c.Request.Context())
	if apiError != nil {
		zap.L().Error("service.GetCommunityList(c.Request.Context()) 错误",
			zap.Error(apiError),
		)
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, list)
}

func CommunityDetailHandler(c *gin.Context) {
	_s := c.Param("id")

	//string 转为 int32
	_t, err := strconv.Atoi(_s)
	if err != nil {
		zap.L().Error("strconv.Atoi(_s) 错误", zap.Error(err))
		ResponseErrorWithCode(c, code.InvalidParam)
		return
	}

	communityID := int32(_t)

	community, apiError := service.GetCommunityDetail(c.Request.Context(), communityID)
	if apiError != nil {
		zap.L().Error("service.GetCommunityDetail(c.Request.Context(), communityID) 错误",
			zap.Error(apiError),
		)
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, community)
	return
}
