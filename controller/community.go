package controller

import (
	"GinTalk/pkg/code"
	"GinTalk/service"
	"github.com/gin-gonic/gin"
)

func CommunityHandler(c *gin.Context) {
	list, apiError := service.GetCommunityList(c.Request.Context())
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, list)
}

func CommunityDetailHandler(c *gin.Context) {
	communityID, exist := c.Get("id")
	if !exist {
		ResponseErrorWithCode(c, code.InvalidParam)
	}
	community, apiError := service.GetCommunityDetail(c.Request.Context(), communityID.(uint))
	if apiError != nil {
		ResponseErrorWithApiError(c, apiError)
		return
	}
	ResponseSuccess(c, community)
}
